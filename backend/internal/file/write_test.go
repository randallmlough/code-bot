package file

import (
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
	// Expected contents of the file
	expectedContents := "Hello, Test!"
	// Path where the file will be created
	filePathName := "testfile.txt"

	// Call the function with the test data
	err := CreateFile(expectedContents, filePathName)
	if err != nil {
		t.Fatalf("Failed to create the file: %v", err)
	}

	// Read the contents of the file
	fileContents, err := os.ReadFile(filePathName)
	if err != nil {
		t.Fatalf("Failed to read the file: %v", err)
	}

	// Verify the contents of the file
	if string(fileContents) != expectedContents {
		t.Fatalf("Unexpected file content. Expected %s but got %s", expectedContents, string(fileContents))
	}

	// Remove the file after testing
	err = os.Remove(filePathName)
	if err != nil {
		t.Fatalf("Failed to remove the file: %v", err)
	}
}
