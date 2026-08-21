package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
)

type shortcut struct {
	ID     string `json:"ID"`
	Key    int    `json:"Key"`
	Name   string `json:"Name"`
	Type   string `json:"Type"`
	Target string `json:"Target"`
}

func Add(key int, name string, pathType string, path string) error {
	entry := shortcut{
		ID:     rand.Text(),
		Key:    key,
		Name:   name,
		Type:   pathType,
		Target: path,
	}
	var shortcuts []shortcut
	file, err := os.OpenFile("shortcuts.json", os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		return err
	}
	shortcuts = append(shortcuts, entry)
	err = json.NewEncoder(file).Encode(&shortcuts)
	if err != nil {
		return err
	}
	return nil
}

func list() error {
	var shortcuts []shortcut
	file, err := os.Open("shortcuts.json")
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		return err
	}
	for _, entry := range shortcuts {
		fmt.Printf("%s - %d - %s - %s - %s", entry.ID, entry.Key, entry.Name, entry.Type, entry.Target)
		fmt.Println()
	}
	return nil
}

func Remove(key int) error {
	var shortcuts []shortcut
	file, err := os.OpenFile("shortcuts.json", os.O_RDWR, 0644)

	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		return err
	}

	for i, entry := range shortcuts {
		if key == entry.Key {
			shortcuts = append(shortcuts[:i], shortcuts[i+1:]...)
			break
		}
	}

	err = json.NewEncoder(file).Encode(&shortcuts)
	if err != nil {
		return err
	}
	return nil
}
