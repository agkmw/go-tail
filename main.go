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

// TODO: Add a flag to show line numbers when output
// TODO: Use channels to communicate data between main and other routines

func main() {
    lineFlag := flag.String("n", "", "Number of lines to display")
    byteFlag := flag.String("c", "", "Number of byteFlag to display (overrides -n)")
    flag.Parse()

    args := flag.Args();
    if len(args) < 1 {
        fmt.Println("Usage: tail [OPTIONS] <filename>")
        os.Exit(1)
    }

    if *lineFlag != "" && *byteFlag != "" {
        fmt.Println("go-tail: Error => You can pass only one flag (-n or -c) at a time")
        return
    }

    if *lineFlag == "" && *byteFlag == "" {
        *lineFlag = "10"
    }

    ch := make(chan []string)
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

            if *byteFlag != "" {
                tailBytes(ch, file, byteFlag, path)
                return
            }
            tailLines(ch, file, lineFlag, path)
        }(path)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()
    for t := range ch {
        for _, lod := range t {
            fmt.Println(lod)
        }
    }
}

func tailBytes(ch chan <- []string, file *os.File, flg *string, path string) {
    byteOffset, err := strconv.Atoi(*flg)
    if err != nil {
        fmt.Println("go-tail: Error => Invalid number for -c flag")
        return
    }

    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Println("go-tail: Error reading file: ", err.Error())
        return
    }

    switch {
    case fileInfo.Size() < int64(byteOffset):
        byteOffset = int(fileInfo.Size())
    case !strings.HasPrefix(*flg, "+"):
        if runtime.GOOS == "windows" {
            byteOffset = byteOffset + 1
        }
        _, err = file.Seek(fileInfo.Size() - int64(byteOffset), 0)
        if err != nil {
            fmt.Println("go-tail: Error reading file: ", err.Error())
            return
        }
    }

    bytesToDisplay := make([]byte, byteOffset)
    _, err = file.Read(bytesToDisplay)
    if err != nil {
        fmt.Println("go-tail: Error reading file: ", err.Error())
        return
    }

    title := fmt.Sprintf("\n==>%v<==\n", path)
    data := []string { title, string(bytesToDisplay) }
    ch <- data
    return
}

func tailLines(ch chan <- []string, file *os.File, flg *string, path string) {
    offset, err := strconv.Atoi(*flg)
    if err != nil {
        fmt.Println("go-tail: Error => Invalid number for -n flag")
        return
    }

    lines := readLines(file)
    var linesToDisplay []string

    switch {
    case len(lines) < offset:
        linesToDisplay = lines
    case strings.HasPrefix(*flg, "+"):
        linesToDisplay = lines[offset:]
    default:
        linesToDisplay = lines[len(lines) - offset:]
    }


    title := fmt.Sprintf("\n==>%v<==\n", path)
    data := append([]string{title}, linesToDisplay...)
    ch <- data
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
