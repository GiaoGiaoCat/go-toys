package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// String 三个参数依次是：参数名称、参数的默认值、参数说明
var path = flag.String("path", "", "项目的绝对路径")

// 入口函数
func main() {
	flag.Parse()
	searchDir := fmt.Sprintf("%s", *path)
	fmt.Println("Working on ", searchDir)

	paths := strings.Split(searchDir, "/")
	outputFile := paths[len(paths)-1]

	fileList, _ := getFileList(searchDir)
	for _, file := range fileList {
		processingData(file, outputFile+"_result.txt")
	}
	fmt.Println("Done.")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// func getFileList(searchDir string) (fileList []string, err error) {
func getFileList(searchDir string) (fileList []string, err error) {
	dirList := []string{}
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		dirList = append(dirList, path)
		return nil
	})
	check(err)

	blackFolders := [...]string{
		"/tmp/", "/spec/", "/doc/", "/log/", "/public/", "/fixtures/",
		"/.git/", "/db/", "/profiles/", "/app/docs/", "/lib/tasks/", "/bin/",
		"/.vscode/",
	}

OuterLoop:
	for _, dir := range dirList {
		for _, folderName := range blackFolders {
			if strings.Contains(dir, folderName) {
				continue OuterLoop
			}

			fi, err := os.Stat(dir)
			check(err)

			mode := fi.Mode()
			if mode.IsDir() {
				continue OuterLoop
			}
		}
		fileList = append(fileList, dir)
	}

	return fileList, nil
}
func writeLines(fileName string, lines []string, outputFile string) error {
	// overwrite file if it exists
	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	check(err)
	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	w.WriteString("文件名：" + fileName + "\n")

	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		check(err)
	}

	w.WriteString("\n")
	// flush outstanding data
	return w.Flush()
}

func processingData(inputFile, outputFile string) {
	file, err := os.Open(inputFile)
	check(err)
	defer file.Close()

	blackLines := [...]string{
		"api :", "error :", "param :", "desc '", "// ", "# ",
	}

	scanner := bufio.NewScanner(file)
	var tmpSlice []string
OuterLoop:
	for scanner.Scan() {
		str := scanner.Text()

		for _, lineStr := range blackLines {
			if strings.Contains(str, lineStr) {
				continue OuterLoop
			}
		}
		if strings.Contains(str, "圈") {
			str = strings.TrimSpace(str)
			tmpSlice = append(tmpSlice, str)
		}
	}
	// fmt.Println(tmpSlice)
	if len(tmpSlice) == 0 {
		return
	}
	err = writeLines(inputFile, tmpSlice, outputFile)
	check(err)
	err = scanner.Err()
	check(err)
}
