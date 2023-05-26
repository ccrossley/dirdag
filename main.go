package main

import (
	"fmt"
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

func printDir(path string, prefix string, depth, maxDepth int) error {
	// Get information about the path
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	// Get the base name of the path
	name := filepath.Base(path)

	// Skip hidden files or directories
	if strings.HasPrefix(name, ".") {
		return nil
	}

	// Check if entry is a symlink
	symlink := ""
	if (info.Mode() & os.ModeSymlink) != 0 {
		symlink = " (symlink)"
		// If it's a symlink, get the information about the path it points to
		info, err = os.Stat(path)
		if err != nil {
			return err
		}
	}

	// Print only directories and .fish files
	if info.IsDir() || filepath.Ext(name) == ".fish" {
		fmt.Println(prefix + name + symlink)
	}

	// If it's a directory and we haven't reached max depth, recurse further
	if info.IsDir() && depth < maxDepth {
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for i, file := range files {
			newPath := filepath.Join(path, file.Name())
			newPrefix := indent
			if i == len(files)-1 {
				newPrefix = lastIndent
			}
			entryPrefix := prefix + newPrefix
			newPrefix = prefix + newPrefix
			if i == len(files)-1 {
				entryPrefix = prefix + lastPrefix
			} else {
				entryPrefix = prefix + prefix
			}
			err = printDir(newPath, entryPrefix, depth+1, maxDepth)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
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
	fmt.Println(root)
	err := printDir(root, prefix, 1, maxDepth)
	if err != nil {
		fmt.Println("Error printing directory:", err)
		os.Exit(1)
	}
	fmt.Println("Diagram generation completed.")
}
