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
		str := scanner.Text()
		if strings.Contains(str, "圈") {
			tmpSlice = append(tmpSlice, str)
		}
	}

	err = writeLines(tmpSlice, outputFile)
	check(err)
	err = scanner.Err()
	check(err)
}

func main() {
	processingData("data.txt", "file.txt")
}
