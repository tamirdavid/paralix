package paralixutils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"paralix/lib/logger"
	"regexp"
	"sort"
	"strings"
)

func CheckAllSubstringsExists(substrings []string, str string) bool {
	// check that all substrings in 'substrings' slice are exists in str
	for _, s := range substrings {
		if !strings.Contains(str, s) {
			return false
		}
	}
	return true
}

func GetMatchedRegexOccurencesFromString(regex string, input string) []string {
	re := regexp.MustCompile(regex)
	// Find all matches of the regular expression in the input string
	matches := re.FindAllStringSubmatch(input, -1)

	// Extract the matched values and add them to a slice of strings
	var result []string
	for _, match := range matches {
		if len((match)) > 1 {
			result = append(result, match[1])
		}
	}
	return result
}

func IsStringInSlice(s []string, target string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, target)
	return i < len(s) && s[i] == target
}

func ReadLinesFromFileReturnSliceOfLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func GetValuesBetweenDelimiters(str string, openChar string, closeChar string, splitChar string) ([]string, error) {
	openCharIndex := strings.Index(str, openChar)
	closeCharIndex := strings.Index(str, closeChar)
	if openCharIndex == -1 || closeCharIndex == -1 {
		return nil, fmt.Errorf("values should be in the format of %sVALUE1%s%sVALUE2%s%sVALUE3%s", openChar, closeChar, splitChar, splitChar, splitChar, closeChar)
	}
	// extract the values between the braces
	valuesStr := str[openCharIndex+len(openChar) : closeCharIndex]
	// split the values string into a slice
	values := strings.Split(valuesStr, splitChar)
	// trim any leading or trailing spaces from the values
	for i, val := range values {
		values[i] = strings.TrimSpace(val)
	}
	return values, nil
}

func RunCmdAndWaitForItToFinish(cmd *exec.Cmd) error {
	logger.Log.Info("Executing command:" + cmd.String())
	if ExecutionErr := cmd.Start(); ExecutionErr != nil {
		fmt.Fprintf(os.Stderr, "Error starting command: %v\n", ExecutionErr)
		return ExecutionErr
	}
	if waitingErr := cmd.Wait(); waitingErr != nil {
		logger.Log.Errorf("Error waiting for command to complete: %v\n", waitingErr)
		return waitingErr
	}
	return nil
}
