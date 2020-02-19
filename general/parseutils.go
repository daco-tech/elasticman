package general

import (
	"fmt"
	"log"
)

// HasPrefix function checks if the prefix is contained in the string.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

// ### AskForConfirmation(), PosString(), ContainsString() - THANKS TO: https://gist.github.com/albrow/5882501

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if ContainsString(okayResponses, response) {
		return true
	} else if ContainsString(nokayResponses, response) {
		return false
	} else {
		log.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func PosString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func ContainsString(slice []string, element string) bool {
	return !(PosString(slice, element) == -1)
}
