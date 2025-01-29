package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func main() {
    nFlag := flag.String("n", "", "Number of lines to display")
    cFlag := flag.String("c", "", "Number of cFlag to display (overrides -n)")
    flag.Parse()

    args := flag.Args();
    if len(args) < 1 {
        fmt.Println("Usage: tail [OPTIONS] <filename>")
        os.Exit(1)
    }

    var wg sync.WaitGroup
    for _, path := range args {
        wg.Add(1)
        go func(path string) {
            defer wg.Done()

            file, err := os.Open(path)
            if err != nil {
                fmt.Println("go-tail: Error reading the file: ", err.Error())
                return
            }
            defer file.Close()

            file = removeBOM(file)

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
                    fmt.Println("go-tail: Error reading file: ", err.Error())
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
                        fmt.Println("go-tail: Error reading file: ", err.Error())
                        return
                    }
                }

                bytesToDisplay := make([]byte, charOffset)
                _, err = file.Read(bytesToDisplay)
                if err != nil {
                    fmt.Println("go-tail: Error reading file: ", err.Error())
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

            fmt.Printf("==>%v<==\n", path)
            for _, line := range linesToDisplay {
                fmt.Println(line)
            }
        }(path)
    }
    wg.Wait()
}

func readLines(file *os.File) []string {
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines
}

func removeBOM(file *os.File) *os.File {
    bom := []byte{0xEF, 0xBB, 0xBF}

    buf := make([]byte, 3)
    _, err := file.Read(buf)
    if err != nil {
        fmt.Println("go-tail: Error reading the file: ", err.Error())
        return file
    }

    var startPos int64 = 0
    if bytes.Equal(buf, bom) {
        startPos = 3
    }

    _, err = file.Seek(startPos, 0)
    if err != nil {
        fmt.Println("go-tail: Error reading file: ", err.Error())
    }
    return file
}
