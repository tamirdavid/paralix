/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"paralix/lib/logger"
	osutils "paralix/lib/osUtils"
	paralixutils "paralix/lib/paralixUtils"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Run a command with N args in parallel.",
	Long:  `Run a command with N args in parallel.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputValidationError := validateCommandInput()
		if inputValidationError != nil {
			return inputValidationError
		}
		outputResourcesError := handleOutputfile()
		if outputResourcesError != nil {
			return outputResourcesError
		}
		executeErr := executeParallel()
		if executeErr != nil {
			return executeErr
		}
		writeResultsErr := writeResultstoFile()
		if writeResultsErr != nil {
			return writeResultsErr
		}
		return nil
	},
}

var command string
var placeholders string
var filepathInput string
var outputfile string
var outputfilesDir string = "/tmp/paralix_output/"

func init() {
	rootCmd.AddCommand(commandCmd)
	commandCmd.Flags().StringVarP(&command, "execute", "e", "", "Command to execute with placeholders (<KEY>) [Example: --execute 'echo <WHAT_SHOULD_ECHO>']")
	commandCmd.Flags().StringVarP(&placeholders, "placeholder", "p", "", "Placeholders in the format of 'KEY={VALUE1,VALUE2,VALUE3}' [Example -p 'WHAT_SHOULD_ECHO={HELLO,WORLD]'")
	commandCmd.Flags().StringVarP(&filepathInput, "inputfile", "f", "", "File that contain the inputs to run, each input in a new line' [Example -f 'customers']")
	commandCmd.Flags().StringVarP(&outputfile, "output", "o", "", "Output file that the results for the command will be written in")
	commandCmd.MarkFlagRequired("output")
	commandCmd.MarkFlagRequired("execute")
}

func writeResultstoFile() error {
	files, err := filepath.Glob(filepath.Join(outputfilesDir, "*"))
	if err != nil {
		return err
	}
	// Concatenate the files and write the result to the output file
	output, creationFileError := osutils.CreateFile(outputfile)
	if creationFileError != nil {
		return creationFileError
	}
	defer output.Close()

	readWriteError := osutils.ReadFilesContentAndWriteToOneFile(output, files)
	if readWriteError != nil {
		return readWriteError
	}
	osutils.PrintFileContent(outputfile)
	osutils.RemoveDirectory(outputfilesDir)
	return nil
}

func handleOutputfile() error {
	_, creationFileError := osutils.CreateFile(outputfile)
	if creationFileError != nil {
		return creationFileError
	}
	dirCreationError := osutils.MakeDirectoryDeleteIfExists(outputfilesDir)
	if dirCreationError != nil {
		return dirCreationError
	}
	return nil
}

func checkIfbothPlaceholdersMethodsUsed() {
	// exit if user passed placeholders and filepath
	if placeholders != "" && filepathInput != "" {
		logger.Log.Error("You can't use both --placeholder [-p] and --inputfile [-f]")
		os.Exit(1)
	}
}

func validateCommandInput() error {
	commandPlaceholders := paralixutils.GetMatchedRegexOccurencesFromString("<(.*?)>", command)
	checkIfbothPlaceholdersMethodsUsed()
	if placeholders != "" {
		if placeHolderError := validatePlaceholderInput(commandPlaceholders); placeHolderError != nil {
			return placeHolderError
		}
		return nil
	} else if filepathInput != "" {
		if placeHolderError := validatePlaceholderFileInput(commandPlaceholders); placeHolderError != nil {
			return placeHolderError
		}
	}
	return nil
}
func validatePlaceholderFileInput(commandPlaceholders []string) error {
	fileName := filepath.Base(filepathInput)
	isExists := paralixutils.IsStringInSlice(commandPlaceholders, fileName)
	if !isExists {
		return errors.New(fmt.Sprintf("<%s> is missing in the command", fileName))
	}
	return nil
}

func validatePlaceholderInput(commandPlaceholders []string) error {
	if !paralixutils.CheckAllSubstringsExists([]string{"=", "{", "}"}, placeholders) {
		return errors.New("Placeholder should be in the format of -p KEY={VALUE1,VALUE2}")
	}
	key := getKeyFromString(placeholders)
	isExists := paralixutils.IsStringInSlice(commandPlaceholders, key)
	if !isExists {
		err := fmt.Sprintf("<%s> is missing in the command", key)
		return errors.New(err)
	}

	// check all command placeholders are passed through -p
	placeholdersKeys := extractKeys(placeholders)
	for _, str := range commandPlaceholders {
		isExists := paralixutils.IsStringInSlice(placeholdersKeys, str)
		if !isExists {
			err := fmt.Sprintf("<%s> has not passed using -p %s=value", str, str)
			return errors.New(err)
		}
	}
	return nil
}

func getKeyFromString(placeholder string) string {
	parts := strings.Split(placeholder, "=")
	if len(parts) != 2 {
		// The input string is not in the correct format
		return ""
	}
	return parts[0]
}

func extractKeys(s string) []string {
	var keys []string
	pairs := strings.Fields(s)
	for _, pair := range pairs {
		keyVal := strings.Split(pair, "=")
		if len(keyVal) != 2 {
			continue
		}
		key := keyVal[0]
		if !paralixutils.IsStringInSlice(keys, key) {
			keys = append(keys, key)
		}
	}
	return keys
}

func valuesToSlice(str string) ([]string, error) {
	// find the opening and closing braces
	if placeholders != "" {
		values, err := paralixutils.GetValuesBetweenDelimiters(str, "{", "}", ",")
		if err != nil {
			return nil, err
		}
		return values, nil
	}
	if filepathInput != "" {
		values, err := paralixutils.ReadLinesFromFileReturnSliceOfLines(filepathInput)
		if err != nil {
			return nil, err
		}
		return values, nil
	}
	return nil, nil
}

func getKeyValuesBasedOnPlaceholderInsertingMethod() (string, []string, error) {
	var toSend string
	var key string
	if placeholders != "" {
		key = getKeyFromString(placeholders)
		toSend = placeholders
	}
	if filepathInput != "" {
		key = filepathInput
		toSend = filepathInput
	}
	values, err := valuesToSlice(toSend)
	return filepath.Base(key), values, err
}

func executeParallel() error {
	key, values, err := getKeyValuesBasedOnPlaceholderInsertingMethod()
	if err != nil {
		return err
	}
	ch := make(chan string)
	for _, placeholder := range values {
		go func(placeholder string) {
			cmdStr := strings.Replace(command, "<"+key+">", placeholder, -1)
			cmd := exec.Command("bash", "-c", cmdStr+" | tee "+outputfilesDir+placeholder)
			// Run command in paralllel report failure to channel if failed [execute/wait]
			executionErr := paralixutils.RunCmdAndWaitForItToFinish(cmd)
			if executionErr != nil {
				ch <- "error"
			}
			ch <- "success"
		}(string(placeholder))
	}
	// wait for all the goroutines to complete
	for range values {
		<-ch
	}
	return nil
}
