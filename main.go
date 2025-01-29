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
    if err != nil {
        fmt.Println("Error reading the file: ", err.Error())
        return
    }
    defer file.Close()

    if *nFlag != "" && *cFlag != "" {
        fmt.Println("Error: You can pass only one flag (-n or -c) at a time")
        return
    }

    if *nFlag == "" && *cFlag == "" {
        *nFlag = "10"
    }

    if *cFlag != "" {
        charOffset, err := strconv.Atoi(*cFlag)
        if err != nil {
            fmt.Println("Error: Invalid number for -c flag")
            return
        }

        fileInfo, err := file.Stat()
        if err != nil {
            fmt.Println("Error reading file: ", err.Error())
            return
        }

        switch {
        case fileInfo.Size() < int64(charOffset):
            charOffset = int(fileInfo.Size())
        case !strings.HasPrefix(*cFlag, "+"):
            if runtime.GOOS == "windows" {
                charOffset = charOffset + 2
            }
            _, err = file.Seek(fileInfo.Size() - int64(charOffset), 0)
            if err != nil {
                fmt.Println("Error reading file: ", err.Error())
                return
            }
        }

        bytesToDisplay := make([]byte, charOffset)
        _, err = file.Read(bytesToDisplay)
        if err != nil {
            fmt.Println("Error reading file: ", err.Error())
            return
        }
        fmt.Println(string(bytesToDisplay))
        return
    }


    offset, err := strconv.Atoi(*nFlag)
    if err != nil {
        fmt.Println("Error: Invalid number for -n flag")
        return
    }

    lines := readLines(file)
    var linesToDisplay []string

    switch {
    case len(lines) < offset:
        linesToDisplay = lines
    case strings.HasPrefix(*nFlag, "+"):
        linesToDisplay = lines[offset:]
    default:
        linesToDisplay = lines[len(lines) - offset:]
    }

    for _, line := range linesToDisplay {
        fmt.Println(line)
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
