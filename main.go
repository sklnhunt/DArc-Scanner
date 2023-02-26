package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Rules struct {
	Rules []string `json:"rules"`
}


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <dir> <targetDir>")
		os.Exit(1)
	}

	// Get the directory name and target directory name from command-line arguments
	dirName := os.Args[1]
	targetDir := os.Args[2]

	// Display status of all files
	displayFilesStatus(dirName)

	// Scan the directory and its subdirectories for code files
	codeFiles := scanDir(dirName)

	// Read the target strings from the target directory
	targetStrings, err := readTargetStrings(targetDir)
	if err != nil {
		fmt.Printf("Error reading target strings: %v\n", err)
		os.Exit(1)
	}

	// Search the code files for the target strings
	for _, file := range codeFiles {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			for _, target := range targetStrings {
				if strings.Contains(scanner.Text(), target) {
					fmt.Printf("Found target string '%s' in file %s\n", target, file)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error scanning file %s: %v\n", file, err)
		}
	}
}

func displayFilesStatus(rootDirToScan string) {

	countsOfExtension := make(map[string]int)
	totalLinesInFiles := 0

	filepath.Walk(rootDirToScan, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			extensionsOfFile := filepath.Ext(info.Name())

			if extensionsOfFile != "" {

				countsOfExtension[extensionsOfFile]++

				// Count the lines of code in the file
				file, err := os.Open(path)
				if err == nil {
					defer file.Close()

					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						totalLinesInFiles++
					}
				}
			}
		}
		return nil
	})

	for ext, count := range countsOfExtension {
		fmt.Printf("Number of file/files containing %s extension: %d\n", ext, count)
	}
	fmt.Println("-------------------------------------------")
	fmt.Printf("Total lines of code: %d\n", totalLinesInFiles)
	fmt.Println("-------------------------------------------")

}

func scanDir(dirName string) []string {
	// Get a list of all code files in the directory and its subdirectories
	var files []string
	filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", dirName, err)
			return err
		}

		if !info.IsDir() && isCodeFile(path) {
			files = append(files, path)
		}

		return nil
	})
	return files
}

func isCodeFile(fileName string) bool {
	// Check if the file has a code extension
	ext := filepath.Ext(fileName)
	switch ext {
	case ".java", ".py", ".html", ".rb":
		return true
	}
	return false
}

func readTargetStrings(targetDir string) ([]string, error) {
	// Read the target strings from the target directory
	var targetStrings []string

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Open the JSON file and read its contents
		jsonFile, err := os.Open(filepath.Join(targetDir, file.Name()))
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Parse the JSON file and extract the target strings
		var rules Rules
		err = json.Unmarshal(byteValue, &rules)
		if err != nil {
			return nil, err
		}

		targetStrings = append(targetStrings, rules.Rules...)
	}

	return targetStrings, nil
}