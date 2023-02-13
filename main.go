package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func displayFilesStatus(rootDirToScan string) map[string]int {

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

    return countsOfExtension
}

func main() {
	rootDirToScan := os.Args[1]

	totalExtensions := displayFilesStatus(rootDirToScan)
    //need to import json files based on totalExtensions Map

	extensionWithRules := map[string][]string{
		".java": {"printStackTrace", "Readline"},
		".py":   {"string3", "string4"},
		".c":    {"string5", "string6"},
		".cpp":  {"string7", "string8"},
		".h":    {"string9", "string10"},
		".go":   {"string11", "string12"},
		".html": {"string13", "string14"},
		".js":   {"string15", "string16"},
		".css":  {"string17", "string18"},
		".xml":  {"string19", "string20"},
	}

	filepath.Walk(rootDirToScan, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			stringsFromRules, err := extensionWithRules[filepath.Ext(info.Name())]
			if !err {
				return nil
			}
			fileToScan, _ := os.Open(path)
			defer fileToScan.Close()
			scanner := bufio.NewScanner(fileToScan)
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				for _, stringToScan := range stringsFromRules {
					if strings.Contains(line, stringToScan) {
						fmt.Printf("[+] Found '%s' in %s at line %d:=> %s\n", stringToScan, path, lineNumber, line)
						break
					}
				}
				lineNumber++
			}
		}
		return nil
	})
}
