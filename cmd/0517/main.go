package main

import (
	"bufio"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeLines(lines []string, outputFile string) error {
	// overwrite file if it exists
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	for _, line := range lines {
		if len(line) < 4 {
			break
		}
		_, err := w.WriteString(line + "\n")
		check(err)
	}

	// flush outstanding data
	return w.Flush()
}

func processingData(inputFile, outputFile string) {
	file, err := os.Open(inputFile)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tmpSlice []string
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), " ")
		tmpSlice = append(tmpSlice, strs...)
	}

	strSlice := uniqueStr(tmpSlice)
	err = writeLines(strSlice, outputFile)
	check(err)
	err = scanner.Err()
	check(err)
}

func uniqueStr(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	processingData("data.txt", "file.txt")
}
