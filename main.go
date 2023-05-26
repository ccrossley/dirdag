package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("Provide directory to diagram")
	}
	root := args[0]
	if !strings.HasSuffix(root, string(os.PathSeparator)) {
		root += string(os.PathSeparator)
	}
	f, err := os.Create("output.mmd")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("graph TB\n")
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) == ".fish" {
			node := strings.ReplaceAll(path, root, "")
			node = strings.ReplaceAll(node, string(os.PathSeparator), "_")
			label := strings.ReplaceAll(path, root, "")
			if d.IsDir() {
				label += "/"
			}
			f.WriteString(node + "[\"" + label + "\"]\n")
			if d.IsDir() || (d.Type()&fs.ModeSymlink != 0) {
				link := filepath.Dir(path)
				link = strings.ReplaceAll(link, root, "")
				link = strings.ReplaceAll(link, string(os.PathSeparator), "_")
				f.WriteString(link + "-->" + node + "\n")
			}
		}
		return nil
	})
}
