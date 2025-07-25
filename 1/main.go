package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := displayDirTree(out, path, printFiles, "")
	if err != nil {
		return fmt.Errorf("error occured while printing dirs: %v", err)
	}

	return nil
}

func displayDirTree(out io.Writer, path string, printFiles bool, margin string) error {
	// открываем новую папку или файл
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open file with path '%s': %v", path, err)
	}
	defer file.Close()

	// получаем список всех вложенных папок и файлов
	files, err := file.ReadDir(-1)
	if err != nil {
		return fmt.Errorf("can't read directory: %v", err)
	}

	if !printFiles {
		files = deleteFiles(files)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// проходимся по всем папкам и файлам
	for i, file := range files {
		var fmtFileSize string

		fileName, fileSize, isDir, err := getFileInfo(file)
		if err != nil {
			fmt.Fprintf(out, "can't get file info: %v\n", err)
			continue
		}

		if !isDir {
			fmtFileSize = convertFileSize(fileSize)
		}

		prefix := "├───"
		currentMargin := margin + "│\t"
		if i == len(files)-1 {
			prefix = "└───"
			currentMargin = margin + "\t"
		}

		fmt.Fprintf(out, "%s%s%s%s\n", margin, prefix, fileName, fmtFileSize)

		if file.IsDir() {
			err := displayDirTree(out, filepath.Join(path, file.Name()), printFiles, currentMargin)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteFiles(files []os.DirEntry) []os.DirEntry {
	var result []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file)
		}
	}

	return result
}

func getFileInfo(file os.DirEntry) (string, int64, bool, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return "", 0, false, fmt.Errorf("can't get file info: %v", err)
	}

	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()
	isDir := fileInfo.IsDir()

	return fileName, fileSize, isDir, nil
}

func convertFileSize(fileSize int64) string {
	if fileSize == 0 {
		return " (empty)"
	}
	return fmt.Sprintf(" (%db)", fileSize)
}
