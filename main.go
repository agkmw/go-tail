package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
    nFlag := flag.String("n", "", "Number of lines to display")
    bytes := flag.String("c", "", "Number of bytes to display (overrides -n)")

    flag.Parse()

    args := os.Args;
    if len(args) < 2 {
        fmt.Println("Usage: tail [OPTIONS] <filename>")
        os.Exit(1)
    }

    path := args[len(args) - 1]

    file, err := os.Open(path)
    defer file.Close()

    if err != nil {
        fmt.Println("An error occured while reading the file: ", err.Error())
        os.Exit(1)
    }

    if *nFlag != "" && *bytes != "" {
        fmt.Println("Error: You can pass only one flag (-n or -c) at a time")
    }

    if *nFlag == "" && *bytes == "" {
        *nFlag = "10"
    }

    contents := readFileContents(file)

    if *nFlag != "" {
        lines, err := strconv.Atoi(*nFlag)
        if err != nil {
            fmt.Println("Error: Invalid number for -n flag")
        }

        if strings.HasPrefix(*nFlag, "+"){
            for _, line := range contents[lines:] {
                fmt.Println(line)
            }
        }

        for _, line := range contents[len(contents) - lines:] {
            fmt.Println(line)
        }
    }

}

func readFileContents(file *os.File) []string {
    var contents []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        contents = append(contents, scanner.Text())
    }
    return contents
}
