package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func countLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	var numLines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && !strings.HasPrefix(line, "//") {
			numLines++
		}
	}
	if err = scanner.Err(); err != nil {
		return 0, err
	}

	return numLines, nil
}

func countLinesInDirectory(directoryPath, fileExtension string) (int, error) {
	var numLines int
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), fileExtension) {
			linesInFile, err := countLinesInFile(path)
			if err != nil {
				return err
			}
			numLines += linesInFile
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return numLines, nil
}

//func main() {
//	totalLines, err := countLinesInDirectory(".", ".go")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	fmt.Printf("Total number of lines in *.go files: %d\n", totalLines)
//}
