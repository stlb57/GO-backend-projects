package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// //filevault list ____

// os

// 1. use os.ReadDir(path)
// 2. returns []DirEntry
// 3. DirEntry is an interface in io/fs
//    -- .Name() string
//    -- .Type() FileMode - filemode represents a file's mode and permission bits
//    -- .Info() (FileInfo, error)
// 4. IsDir() to check if the current file is a folder?
// 5. path/filepath??
// 6.

func list(path string) error {
	list, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	fmt.Println("---------------Files-----------------")
	fmt.Println("\n\n")
	file_count := 1
	folder_count := 1
	for _, entry := range list {
		if !entry.IsDir() {
			fmt.Printf("%d . %s", file_count, entry.Name())
			fmt.Println()
			file_count++
		}
	}
	fmt.Println("\n\n\n\n")
	fmt.Println("---------------Directories-----------------")
	fmt.Println("\n\n")
	for _, entry := range list {
		if entry.IsDir() {
			fmt.Printf("%d . %s", folder_count, entry.Name())
			fmt.Println()
			folder_count++
		}
	}
	return nil
}

func fileInfo(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	fmt.Printf("Name : %s", info.Name())
	fmt.Println()
	fmt.Printf("Size : %d", info.Size())
	fmt.Println()
	fmt.Printf("ModTime : %v", info.ModTime())
	fmt.Println()
	fmt.Printf("IsDir : %v", info.IsDir())
	fmt.Println()
	return nil

}

func search(path string, key string) {
	filepath.WalkDir(
		path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && d.Name() == key {
				fmt.Printf("Found at location : %s", path)
			}
			return nil
		},
	)
}

func delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// Copy function

// 1. os.Open file
// 2. create new file at destination
// 3. io.Reader from file to be copied, copy contents
// 4. io.Writer to the destined file
// 5. close file
// 6. print copied

func Copy(source string, destination string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}

func Move(source string, destination string) error {
	err := os.Rename(source, destination)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := fileInfo(`C:\Users\saart\Downloads\tiffin-connect-app.md`); err != nil {
		fmt.Println(err)
		return
	}

}
