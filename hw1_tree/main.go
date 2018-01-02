package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const separator = string(os.PathSeparator)

func FilterDirs(vs []os.FileInfo) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if v.IsDir() {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func tree(w io.Writer, path string, showFiles bool, padding string) error {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	if !showFiles {
		list = FilterDirs(list)
	}

	var count int = 0
	listLength := len(list)

	for _, file := range list {
		count++
		if file.IsDir() {
			if count == listLength {
				fmt.Fprintf(w, "%s%s%s\n", padding, "└───", file.Name())
				tree(w, path+separator+file.Name(), showFiles, padding+"	")
			} else {
				fmt.Fprintf(w, "%s%s%s\n", padding, "├───", file.Name())
				tree(w, path+separator+file.Name(), showFiles, padding+"│	")
			}
		} else {
			var fileSize string
			if file.Size() == 0 {
				fileSize = "(empty)"
			} else {
				fileSize = "(" + fmt.Sprintf("%d", file.Size()) + "b)"
			}

			if count == listLength {
				fmt.Fprintf(w, "%s%s%s %s\n", padding, "└───", file.Name(), fileSize)
			} else {
				fmt.Fprintf(w, "%s%s%s %s\n", padding, "├───", file.Name(), fileSize)
			}
		}
	}
	return nil
}

func dirTree(w io.Writer, path string, printFiles bool) error {
	err := tree(w, path, printFiles, "")
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
