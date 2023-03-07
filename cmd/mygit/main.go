package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		opt := os.Args[2]
		switch opt {
		case "-p":
			blobsha := os.Args[3]
			fpath := filepath.Join(".git/objects", blobsha[:2], blobsha[2:])
			f, err := os.Open(fpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening %s: %s\n", fpath, err)
				os.Exit(1)
			}
			a, err := zlib.NewReader(f)
			if err != nil {
				log.Fatal(err)
			}
			defer a.Close()

			b, err := io.ReadAll(a)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(strings.Split(string(b), "\x00")[1])

		default:
			fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
