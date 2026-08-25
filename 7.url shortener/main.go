package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
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

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(shortcuts)
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

func Edit(key int, field string, value string) error {
	var shortcuts []shortcut

	file, err := os.Open("shortcuts.json")
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&shortcuts)
	file.Close()

	if err != nil {
		return err
	}

	for i := range shortcuts {
		if shortcuts[i].Key == key {
			switch field {
			case "name":
				shortcuts[i].Name = value
			case "type":
				shortcuts[i].Type = value
			case "target":
				shortcuts[i].Target = value
			case "key":
				return fmt.Errorf("key cannot be edited")
			default:
				return fmt.Errorf("invalid field")
			}
			break
		}
	}

	file, err = os.Create("shortcuts.json")
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(shortcuts)
	if err != nil {
		return err
	}

	return nil
}

func Run(identifier string) error {
	var shortcuts []shortcut

	file, err := os.Open("shortcuts.json")
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&shortcuts)
	file.Close()

	if err != nil {
		return err
	}

	for _, entry := range shortcuts {
		if fmt.Sprint(entry.Key) == identifier || entry.ID == identifier {
			switch entry.Type {
			case "url":
				return exec.Command("cmd", "/c", "start", "", entry.Target).Start()

			case "folder":
				return exec.Command("explorer.exe", entry.Target).Start()

			case "app":
				return exec.Command(entry.Target).Start()

			default:
				return fmt.Errorf("invalid shortcut type")
			}
		}
	}

	return fmt.Errorf("shortcut not found")
}

func buildIndex(shortcuts []shortcut) map[string]shortcut {
	index := make(map[string]shortcut)

	for _, entry := range shortcuts {
		index[entry.ID] = entry
	}

	return index
}

func validateShortcut(key int, name string, shortcutType string, target string) error {
	var errs []string
	if key <= 0 {
		errs = append(errs, "Key should be a positive number")
	}
	if strings.TrimSpace(name) == "" {
		errs = append(errs, "Name cannot be empty string")
	}
	if !(shortcutType == "url" || shortcutType == "folder" || shortcutType == "app") {
		errs = append(errs, "Invalid Shortcut type")
	}
	if strings.TrimSpace(target) == "" {
		errs = append(errs, "Name cannot be empty string")
	}
	if shortcutType == "url" {
		_, err := url.ParseRequestURI(target)
		if err != nil {
			errs = append(errs, "Invalid URL")
		}
	}

	if shortcutType == "folder" {
		info, err := os.Stat(target)
		if err != nil || !info.IsDir() {
			errs = append(errs, "Invalid folder")
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func main() {
	err := Add(1, "youtube", "url", "https://youtube.com")
	if err != nil {
		fmt.Println(err)
	}

	err = list()
	if err != nil {
		fmt.Println(err)
	}

	err = Run("1")
	if err != nil {
		fmt.Println(err)
	}
}
