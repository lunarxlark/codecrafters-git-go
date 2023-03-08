package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {

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
			zr, err := zlib.NewReader(f)
			if err != nil {
				log.Fatal(err)
			}
			defer zr.Close()

			b, _ := io.ReadAll(zr)
			fmt.Print(strings.Split(string(b), "\x00")[1])

		default:
			fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
			os.Exit(1)
		}

	case "hash-object":
		opt := os.Args[2]
		switch opt {
		case "-w":
			fname := os.Args[3]
			fmt.Println(fname)
			digest := sha1.New()
			digest.Write([]byte(fname))
			a := hex.EncodeToString(digest.Sum(nil))
			fmt.Println(a)

		default:
			fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
