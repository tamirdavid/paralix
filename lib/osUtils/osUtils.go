package osUtils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tamirdavid/paralix/lib/logger"
)

func RemoveDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(filePath string) (*os.File, error) {
	output, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func MakeDirectoryDeleteIfExists(dirName string) error {
	// Check if directory already exists
	if _, err := os.Stat(dirName); err == nil {
		// Directory exists, so delete it
		if err := os.RemoveAll(dirName); err != nil {
			return err
		}
	}
	// Create the directory
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return err
	}

	return nil
}

func PrintFileContent(filePath string) error {
	logger.Log.Infof("Printing the content of file %s: \n", filePath)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	fmt.Println(string(content))
	return nil
}

func ReadFilesContentAndWriteToOneFile(output *os.File, files []string) error {
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%s\n%s\n", filepath.Base(file), content)
	}
	return nil
}
