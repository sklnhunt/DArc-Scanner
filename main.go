package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
	// Get the directory name and target directory name using flag
	dirToScan := flag.String("d","","directory to scan code files")
	dirOfRules := flag.String("r","","directory containing rules files")

	flag.Parse()

	if *dirToScan == "" || *dirOfRules == "" {
		fmt.Println("Usage: go run main.go -d <dirToScan> -r <dirWithRules>")
		os.Exit(1)
	}

	// Display status of all files
	displayFilesStatus(*dirToScan)

	// Scan the directory and its subdirectories for code files
	codeFiles := scanDir(*dirToScan)

	// Read the target strings from the target directory
	targetStrings, err := readTargetStrings(*dirOfRules)
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

func scanDir(dirToScan string) []string {
	// Get a list of all code files in the directory and its subdirectories
	var files []string
	filepath.Walk(dirToScan, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", dirToScan, err)
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

func readTargetStrings(dirOfRules string) ([]string, error) {
	// Read the target strings from the target directory
	var targetStrings []string

	files, err := ioutil.ReadDir(dirOfRules)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Open the JSON file and read its contents
		jsonFile, err := os.Open(filepath.Join(dirOfRules, file.Name()))
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