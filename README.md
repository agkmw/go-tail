# Go-Tail

A simple implementation of `tail` command from Linux in Go.
It can be used to display the last few lines or bytes of a file.
It supports reading multiple files concurrently and offer options
to specify the number of lines or bytes to display.

## Features
- Display the last `n` lines or `c` bytes of a file.
- Read multiple files concurrently.

## Installation
1. Clone the repository:
    ```bash
    git clone git@github.com:agkmw/go-tail.git
    cd go-tail
    ```
2. Build the program:
    ```bash
    go build
    ```
 - Build the program with different name (`<filename>.exe` on Windows):
	```bash
	go build -o <filename>
	```
3. Or install the executable to run it anywhere on your system:
	```bash
	go install
	```

## Usage
```bash
./go-tail [OPTIONS] <filename>
```

## Options
`-n <number>` : Display the last `<number>` of lines.
`-c <number>` : Display the last `<number>` of bytes.
`+<number>` : Display from the `<number>th` line or byte onward.

## Examples
- Display the last 10 lines of a file (default):
	```bash
	./go-tail file.txt
	```
- Display the last 20 lines of a file:
	```bash
	./go-tail -n 20 file.txt
	```
- Display the last 50 bytes of a file:
	```bash
	./go-tail -c 50 file.txt
	```
- Display from the 5th line onward:
	```bash
	./go-tail -n +5 file.txt
	```
- Display multiple files:
	```bash
	./go-tail file.txt file2.txt file3.txt
	```
