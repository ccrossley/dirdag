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

	// Check if entry is a symlink
	symlink := ""
	isSymlink := node.Type()&fs.ModeSymlink != 0
	if isSymlink {
		symlink = " (symlink)"
		var err error
		resolvedPath, err := os.Readlink(filepath.Join(path, node.Name()))
		if err != nil {
			return err
		}
		// Trim the trailing slash if it exists
		resolvedPath = strings.TrimSuffix(resolvedPath, "/")
		// Convert resolved path to absolute path
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Join(path, resolvedPath)
		}

		resolvedInfo, err := os.Stat(resolvedPath)
		if err != nil {
			return err
		}

		// If resolved path is a directory and we haven't reached max depth, recurse further
		if resolvedInfo.IsDir() && depth < maxDepth {
			dirEntries, err := os.ReadDir(resolvedPath)
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
				if isLast {
					entryPrefix = prefix + lastPrefix
				}

				newDepth := depth + 1
				err = printDir(resolvedPath, entry, entryPrefix, newDepth, maxDepth)
				if err != nil {
					return err
				}
			}
		}
		return nil // Skip printing and recursing for symlinked files
	}

	// Print only directories and .fish files
	if node.IsDir() || filepath.Ext(node.Name()) == ".fish" {
		fmt.Println(prefix + node.Name() + symlink)
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
			if isLast {
				entryPrefix = prefix + lastPrefix
			}

			newDepth := depth + 1
			err = printDir(newPath, entry, entryPrefix, newDepth, maxDepth)
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
		err = printDir(root, entry, prefix, 1, maxDepth)
		if err != nil {
			fmt.Println("Error printing directory:", err)
			os.Exit(1)
		}
	}
	fmt.Println("Diagram generation completed.")
}
