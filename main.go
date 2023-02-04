package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    rootDirToScan := os.Args[1]
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
                        fmt.Printf("[+] Found '%s' in %s at line %d: %s\n", stringToScan, path, lineNumber, line)
                        break
                    }
                }
                lineNumber++
            }
        }
        return nil
    })
}
