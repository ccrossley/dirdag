package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Node struct {
	Name     string
	Children []*Node
	IsLink   bool
}

var (
	prefix       = "├── "
	indent       = "│   "
	lastPrefix   = "└── "
	lastIndent   = "    "
	defaultDepth = 3
)

func buildTree(root string, maxDepth int) (*Node, error) {
	rootNode := &Node{Name: root}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Skip hidden files or directories
		parts := strings.Split(rel, string(filepath.Separator))
		for _, part := range parts {
			if strings.HasPrefix(part, ".") {
				return nil
			}
		}

		// Check the depth
		depth := len(parts)
		if depth > maxDepth {
			return nil
		}

		// Check the file type
		if info.IsDir() || filepath.Ext(info.Name()) == ".fish" {
			// Traverse the tree to the correct location
			current := rootNode
			for _, part := range parts {
				found := false
				for _, child := range current.Children {
					if child.Name == part {
						current = child
						found = true
						break
					}
				}
				// Add a new node only if depth is within the limit
				if !found && depth <= maxDepth {
					newNode := &Node{Name: part}
					current.Children = append(current.Children, newNode)
					current = newNode
				}
			}
			current.IsLink = info.Mode()&os.ModeSymlink != 0
		}

		return nil
	})

	return rootNode, err
}

func printTree(node *Node, prefix string) {
	fmt.Println(prefix + node.Name)
	newPrefix := prefix + indent
	if node.IsLink {
		newPrefix += "(symlink) "
	}
	for i, child := range node.Children {
		if i == len(node.Children)-1 {
			printTree(child, newPrefix+lastPrefix)
		} else {
			printTree(child, newPrefix+prefix)
		}
	}
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

	rootNode, err := buildTree(root, maxDepth)
	if err != nil {
		fmt.Println("Error building directory tree:", err)
		os.Exit(1)
	}

	printTree(rootNode, "")
	fmt.Println("Diagram generation completed.")
}
