package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Options struct {
	fileName string
	cFlag    *bool
	lFlag    *bool
	wFlag    *bool
}

type Results struct {
	byteCount int
	wordCount int
	lineCount int
}

func main() {
	options := parseArguments()
	results := readFileAndCount(options.fileName)
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

	if len(args) == 0 {
		fmt.Println("Usage: gowc [-c] <filename>")
		os.Exit(1)
	}

	options.fileName = args[len(args)-1]

	options.cFlag = flag.Bool("c", false, "count the number of bytes in the file")
	options.lFlag = flag.Bool("l", false, "count the number of lines in the file")
	options.wFlag = flag.Bool("w", false, "count the number of words in the file")
	flag.Parse()

	// By default, if no flags are provided, then we will count everything, lines, bytes, etc.
	if !*options.cFlag && !*options.lFlag && !*options.wFlag {
		*options.cFlag, *options.lFlag, *options.wFlag = true, true, true
	}

	return
}

func readFileAndCount(fileName string) (results Results) {
	f, err := os.Open(fileName)
	handleError(err)
	defer f.Close()

	buf := make([]byte, 4*1024)

	for {
		n, err := f.Read(buf)

		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				handleError(err)
			}
		}

		if n > 0 {
			results.byteCount += n
		}

		inWord := true
		for i := range n {
			if buf[i] == '\n' {
				results.lineCount++
			}

			if unicode.IsSpace(rune(buf[i])) {
				if inWord {
					results.wordCount++
					inWord = false
				}
			} else {
				inWord = true
			}
		}
	}

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

	fmt.Println(output + options.fileName)
}
