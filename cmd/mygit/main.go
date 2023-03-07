package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
			fpath := fmt.Sprintf(".git/objects/%s/%s", blobsha[0:2], blobsha[2:])
			f, err := os.Open(fpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening .git/objects/%s/%s: %s\n", blobsha[0:2], blobsha[2:], err)
				os.Exit(1)
			}
			a, _ := zlib.NewReader(f)
			io.ReadAll(a)

		default:
			fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
