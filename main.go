package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	prefix     = "├── "
	indent     = "│   "
	lastPrefix = "└── "
	lastIndent = "    "
)

func printDir(path string, node fs.DirEntry, prefix string) error {
	fmt.Println(prefix + node.Name())
	if node.IsDir() {
		newPath := filepath.Join(path, node.Name())
		dirEntries, err := os.ReadDir(newPath)
		if err != nil {
			return err
		}

		for i, entry := range dirEntries {
			isLast := i == len(dirEntries)-1
			if entry.IsDir() || filepath.Ext(entry.Name()) == ".fish" {
				newPrefix := indent
				if isLast {
					newPrefix = lastIndent
				}
				entryPrefix := prefix + newPrefix
				newPrefix = prefix + newPrefix
				if isLast {
					entryPrefix = prefix + lastPrefix
				} else {
					entryPrefix = prefix + prefix
				}
				err := printDir(newPath, entry, entryPrefix)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Provide directory to diagram")
		os.Exit(1)
	}
	root := args[0]
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}
	fmt.Println(root)
	for i, entry := range dirEntries {
		isLast := i == len(dirEntries)-1
		if entry.IsDir() || filepath.Ext(entry.Name()) == ".fish" {
			prefix := prefix
			if isLast {
				prefix = lastPrefix
			}
			err := printDir(root, entry, prefix)
			if err != nil {
				fmt.Println("Error printing directory:", err)
				os.Exit(1)
			}
		}
	}
	fmt.Println("Diagram generation completed.")
}
