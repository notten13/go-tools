package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	fileName string
	cFlag    *bool
	lFlag    *bool
	wFlag    *bool
	mFlag    *bool
}

type Results struct {
	byteCount int
	wordCount int
	lineCount int
	charCount int
}

func main() {
	options := parseArguments()
	results := getCounts(options.fileName)
	printOutput(options, results)
}

func handleError(e error) {
	if e != nil {
		fmt.Println("gowc: " + e.Error())
		os.Exit(1)
	}
}

func parseArguments() (options Options) {
	args := os.Args[1:]

	if len(args) > 0 && args[len(args)-1][0] != '-' {
		options.fileName = args[len(args)-1]
	}

	options.cFlag = flag.Bool("c", false, "count the number of bytes in the file")
	options.lFlag = flag.Bool("l", false, "count the number of lines in the file")
	options.wFlag = flag.Bool("w", false, "count the number of words in the file")
	options.mFlag = flag.Bool("m", false, "count the number of characters in the file")
	flag.Parse()

	// By default, if no flags are provided, then we will count bytes, lines, and words.
	if !*options.cFlag && !*options.lFlag && !*options.wFlag && !*options.mFlag {
		*options.cFlag, *options.lFlag, *options.wFlag = true, true, true
	}

	return
}

func getCounts(fileName string) (results Results) {
	var reader io.Reader

	if fileName == "" {
		reader = os.Stdin
	} else {
		var err error
		reader, err = os.Open(fileName)
		handleError(err)
		defer reader.(*os.File).Close()
	}

	buf := make([]byte, 4*1024)

	var leftover []byte

	for {
		n, err := reader.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				handleError(err)
			}
		}

		results.byteCount += n
		results.lineCount += countLines(buf[:n])
		results.wordCount += countWords(buf[:n])
		results.charCount += countCharacters(buf[:n], &leftover)
	}

	return
}

func countLines(buffer []byte) (lineCount int) {
	for i := range buffer {
		if buffer[i] == '\n' {
			lineCount++
		}
	}

	return
}

func countWords(buffer []byte) (wordCount int) {
	inWord := true
	for i := range buffer {
		if unicode.IsSpace(rune(buffer[i])) {
			if inWord {
				wordCount++
				inWord = false
			}
		} else {
			inWord = true
		}
	}

	return
}

func countCharacters(buffer []byte, leftover *[]byte) (charCount int) {
	i := 0
	data := append(*leftover, buffer...)

	for i < len(data) {
		char, size := utf8.DecodeRune(data[i:])
		if char == utf8.RuneError && size == 1 {
			break
		}
		charCount++
		i += size
	}

	*leftover = data[i:]

	return
}

func printOutput(options Options, results Results) {
	output := "  "

	if *options.lFlag {
		output += strconv.Itoa(results.lineCount) + " "
	}

	if *options.wFlag {
		output += strconv.Itoa(results.wordCount) + " "
	}

	if *options.cFlag {
		output += strconv.Itoa(results.byteCount) + " "
	}

	if *options.mFlag {
		output += strconv.Itoa(results.charCount) + " "
	}

	fmt.Println(output + options.fileName)
}
