package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

var debug bool

// usage prints how to use the program
func usage() {
	fmt.Println("Usage: tde <changed_files> <target_dir>")
}

// checkFileExists checks if the file exists and returns an error if not
func checkFileExists(file string) error {
	_, err := os.Stat(file)
	return err
}

func ExtractPathFromURL(url string) string {
	re := regexp.MustCompile(`\.git/+(.*?)(\?|$)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 && match[1] != "" {
		if debug {
			fmt.Println("[DEBUG] Match found:", match[1])
		}
		return "/" + match[1]
	} else {
		return ""
	}
}

// findModuleSource extracts "source" attributes from module blocks in the HCL body
func findModuleSource(body *hclwrite.Body) []string {
	var moduleSources []string
	for _, block := range body.Blocks() {
		if block.Type() == "module" {
			if attr := block.Body().GetAttribute("source"); attr != nil {
				// Extract source value and clean up surrounding quotes
				source := strings.Trim(string(attr.Expr().BuildTokens(nil).Bytes()), "\"")
				if debug {
					fmt.Println("[DEBUG] Found module source:", source)
				}
				source = strings.TrimPrefix(source, " \"")
				if source == "" || strings.HasPrefix(source, "https://") {
					continue
				}
				if strings.HasPrefix(source, "git@") {
					source = ExtractPathFromURL(source)
				}
				moduleSources = append(moduleSources, source)
			}
		}
	}
	return moduleSources
}

// checkSuffixMatch checks if the end parts of two paths match
func checkSuffixMatch(path1, path2 string) bool {
	parts1 := strings.Split(filepath.Clean(path1), string(filepath.Separator))
	parts2 := strings.Split(filepath.Clean(path2), string(filepath.Separator))

	minLen := len(parts1)
	if len(parts2) < minLen {
		minLen = len(parts2)
	}

	for i := 0; i < minLen; i++ {
		if parts1[len(parts1)-1-i] != parts2[len(parts2)-1-i] {
			return false
		}
	}
	return true
}

// findTFFile walks through a directory and collects .tf files
func findTFFile(dirPath string) ([]string, error) {
	var tfFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" && !strings.HasPrefix(info.Name(), ".") && !strings.Contains(path, ".terraform") {
			tfFiles = append(tfFiles, path)
		}
		return nil
	})

	if debug {
		fmt.Println("[DEBUG] Found .tf files:", tfFiles)
	}
	return tfFiles, err
}

func main() {
	if debug {
		fmt.Println("[DEBUG] Debug mode enabled")
	}

	// Check command-line arguments
	if len(os.Args) < 2 {
		usage()
		return
	}

	changedFiles := os.Args[1]
	if err := checkFileExists(changedFiles); err != nil {
		fmt.Println("File not found:", changedFiles)
		return
	}

	targetDir := "./"
	if len(os.Args) > 2 {
		targetDir = os.Args[2]
	}

	files, err := findTFFile(targetDir)
	if err != nil {
		log.Fatalf("Error finding .tf files: %v", err)
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", file, err)
		}

		parsedFile, diags := hclwrite.ParseConfig(content, file, hcl.InitialPos)
		if diags.HasErrors() {
			log.Fatalf("HCL parse error in file %s: %v", file, diags.Error())
		}

		body := parsedFile.Body()
		sources := findModuleSource(body)
		if len(sources) > 0 {
			for _, source := range sources {
				if checkSuffixMatch(source, filepath.Dir(changedFiles)) {
					fmt.Println("Matching file:", file)
					break
				}
			}
		}
	}
}

func init() {
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
}
