package main

import "os"

func mini_init() error {
	err := os.MkdirAll(".git/refs/heads", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(".git/refs/tags", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(".git/objects", 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(".git/HEAD", []byte("ref: refs/heads/main\n"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if os.Args[1] == "init" {
		err := mini_init()
		if err != nil {
			return
		}
	}
}
