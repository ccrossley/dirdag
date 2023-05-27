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

func printDir(root string, maxDepth int) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore the root directory
		if path == root {
			return nil
		}
		// Get the relative path
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// Skip hidden files and directories
		parts := strings.Split(rel, string(filepath.Separator))
		for _, part := range parts {
			if strings.HasPrefix(part, ".") {
				return filepath.SkipDir
			}
		}
		// Check the depth
		depth := len(parts)
		if depth > maxDepth {
			return filepath.SkipDir
		}
		// Check the file type
		if info.IsDir() || filepath.Ext(info.Name()) == ".fish" {
			// Create the prefix
			prefix := strings.Repeat(indent, depth-1)
			if depth > 1 {
				prefix += prefix
			}
			if info.Mode()&os.ModeSymlink != 0 {
				prefix += "(symlink) "
			}
			fmt.Println(prefix + info.Name())
		}
		return nil
	})
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
	err := printDir(root, maxDepth)
	if err != nil {
		fmt.Println("Error printing directory:", err)
		os.Exit(1)
	}
	fmt.Println("Diagram generation completed.")
}
