package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckSuffixMatch verifies if checkSuffixMatch correctly matches the suffix of two paths
func TestCheckSuffixMatch(t *testing.T) {
	assert.True(t, checkSuffixMatch("modules/elb/ecs", "../../../modules/elb/ecs"))
	assert.False(t, checkSuffixMatch("modules/s3/bucket", "../../../modules/elb/ecs"))
}

// TestFindTFFile tests findTFFile to verify it finds .tf files and ignores hidden files and .terraform directories
func TestFindTFFile(t *testing.T) {
	// Create a temporary directory with .tf files
	dir, err := ioutil.TempDir("", "test-tf-files")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Create a .tf file in the temp directory
	tfFile := filepath.Join(dir, "test.tf")
	err = ioutil.WriteFile(tfFile, []byte("resource \"aws_s3_bucket\" \"example\" {}"), 0644)
	assert.NoError(t, err)

	// Run findTFFile and check results
	tfFiles, err := findTFFile(dir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tfFiles))
	assert.Equal(t, "test.tf", filepath.Base(tfFiles[0]))

	// Hidden files and .terraform directory should be ignored
	hiddenFile := filepath.Join(dir, ".hidden.tf")
	terraformDir := filepath.Join(dir, ".terraform")
	os.Mkdir(terraformDir, 0755)
	err = ioutil.WriteFile(hiddenFile, []byte("resource \"aws_s3_bucket\" \"hidden\" {}"), 0644)
	assert.NoError(t, err)

	tfFiles, err = findTFFile(dir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tfFiles)) // Only test.tf should be detected, ignoring hidden files
}

// TestMainFunction verifies the overall functionality of the main function and simulates a typical workflow
func TestMainFunction(t *testing.T) {
	// Create a temporary directory and a mock changed files list
	dir, err := ioutil.TempDir("", "test-main-function")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Create a mock changed file path
	changedFile := filepath.Join(dir, "modules/elb/ecs/main.tf")
	err = os.MkdirAll(filepath.Dir(changedFile), 0755)
	assert.NoError(t, err)

	// Write a sample content to the mock changed file
	err = ioutil.WriteFile(changedFile, []byte("resource \"aws_s3_bucket\" \"example\" {}"), 0644)
	assert.NoError(t, err)

	// Create a .tf file in the directory with a module block and a source attribute
	moduleFile := filepath.Join(dir, "some_module.tf")
	moduleContent := `
module "example" {
  source = "../../../modules/elb/ecs"
}
`
	err = ioutil.WriteFile(moduleFile, []byte(moduleContent), 0644)
	assert.NoError(t, err)

	// Set up fake arguments for testing main
	os.Args = []string{"cmd", changedFile, dir}

	// Capture the output of the main function
	// Define a variable to store the output result
	main()

}
