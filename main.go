package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
    nFlag := flag.String("n", "", "Number of lines to display")
    cFlag := flag.String("c", "", "Number of cFlag to display (overrides -n)")

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

    if *nFlag != "" && *cFlag != "" {
        fmt.Println("Error: You can pass only one flag (-n or -c) at a time")
    }

    if *nFlag == "" && *cFlag == "" {
        *nFlag = "10"
    }


    if *nFlag != "" {
        contents := readLines(file)
        lines, err := strconv.Atoi(*nFlag)
        if err != nil {
            fmt.Println("Error: Invalid number for -n flag")
        }

        if strings.HasPrefix(*nFlag, "+"){
            for _, line := range contents[lines:] {
                fmt.Println(line)
            }
        } else {
            for _, line := range contents[len(contents) - lines:] {
                fmt.Println(line)
            }
        }
    }

    if *cFlag != "" {
        chars, err := strconv.Atoi(*cFlag)
        if err != nil {
            fmt.Println("Error: Invalid number for -c flag")
        }

        if strings.HasPrefix(*cFlag, "+") {
            contents := make([]byte, chars)
            _, err := file.Read(contents)
            if err != nil {
                fmt.Println("Error reading file: ", err.Error())
            }
            fmt.Println(string(contents))
        } else {
            fileSize, _ := file.Stat()
            fmt.Println(chars)
            useros := runtime.GOOS
            if useros != "windows" {
                _, _ = file.Seek(fileSize.Size() - int64(chars), 0)
            } else {
                _, _ = file.Seek(fileSize.Size() - int64(chars + 2), 0)
            }
            contents := make([]byte, chars)
            _, err := file.Read(contents)
            fmt.Println(len(contents))
            if err != nil {
                fmt.Println("Error reading file: ", err.Error())
            }
            fmt.Println(string(contents))
        }
    }
}

func readLines(file *os.File) []string {
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines
}
