package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	prefix       = "├── "
	indent       = "│   "
	lastPrefix   = "└── "
	lastIndent   = "    "
	defaultDepth = 3
)

func printDir(path string, node fs.DirEntry, prefix string, depth, maxDepth int) error {
	// Skip hidden files or directories
	if strings.HasPrefix(node.Name(), ".") {
		return nil
	}

	// Check if the node is a symlink
	info, err := node.Info()
	if err != nil {
		return err
	}

	// Print only directories and .fish files
	if node.IsDir() || filepath.Ext(node.Name()) == ".fish" {
		if info.Mode()&fs.ModeSymlink != 0 {
			target, err := os.Readlink(filepath.Join(path, node.Name()))
			if err != nil {
				return err
			}

			// Make the target path relative to $GOPATH
			gopath, _ := os.LookupEnv("GOPATH")
			target = strings.Replace(target, gopath, "$GOPATH", 1)

			// Print the symlink with its target
			fmt.Println(prefix + node.Name() + " -> " + target)
		} else {
			fmt.Println(prefix + node.Name())
		}
	}

	// If it's a directory and we haven't reached max depth, recurse further
	if node.IsDir() && depth < maxDepth {
		newPath := filepath.Join(path, node.Name())
		dirEntries, err := os.ReadDir(newPath)
		if err != nil {
			return err
		}

		for i, entry := range dirEntries {
			isLast := i == len(dirEntries)-1
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
			err := printDir(newPath, entry, entryPrefix, depth+1, maxDepth)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("The fug?")
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Provide directory to diagram")
		os.Exit(1)
	}
	root := args[0]
	maxDepth := defaultDepth
	if len(args) > 1 {
		var err error
		maxDepth, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid maximum depth:", err)
			os.Exit(1)
		}
	}
	dirEntries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}
	fmt.Println(root)
	for i, entry := range dirEntries {
		isLast := i == len(dirEntries)-1
		prefix := prefix
		if isLast {
			prefix = lastPrefix
		}

		// Use os.Stat to get information about the entry (os.Stat follows symlinks)
		info, err := os.Stat(filepath.Join(root, entry.Name()))
		if err != nil {
			fmt.Println("Error accessing entry:", err)
			os.Exit(1)
		}

		// Create a fs.DirEntry from the FileInfo obtained from os.Stat
		entry = fs.FileInfoToDirEntry(info)
		err = printDir(root, entry, prefix, 1, maxDepth)
		if err != nil {
			fmt.Println("Error printing directory:", err)
			os.Exit(1)
		}
	}
	fmt.Println("Diagram generation completed.")
}
