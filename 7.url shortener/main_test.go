package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateShortcut(t *testing.T) {
	tests := []struct {
		name         string
		key          int
		shortcutName string
		shortcutType string
		target       string
		wantErr      bool
	}{
		{
			name:         "valid URL",
			key:          1,
			shortcutName: "youtube",
			shortcutType: "url",
			target:       "https://youtube.com",
			wantErr:      false,
		},
		{
			name:         "valid folder",
			key:          2,
			shortcutName: "temp",
			shortcutType: "folder",
			target:       os.TempDir(),
			wantErr:      false,
		},
		{
			name:         "valid app",
			key:          3,
			shortcutName: "notepad",
			shortcutType: "app",
			target:       "notepad.exe",
			wantErr:      false,
		},
		{
			name:         "invalid key",
			key:          0,
			shortcutName: "youtube",
			shortcutType: "url",
			target:       "https://youtube.com",
			wantErr:      true,
		},
		{
			name:         "negative key",
			key:          -1,
			shortcutName: "youtube",
			shortcutType: "url",
			target:       "https://youtube.com",
			wantErr:      true,
		},
		{
			name:         "empty name",
			key:          1,
			shortcutName: "   ",
			shortcutType: "url",
			target:       "https://youtube.com",
			wantErr:      true,
		},
		{
			name:         "invalid type",
			key:          1,
			shortcutName: "youtube",
			shortcutType: "banana",
			target:       "https://youtube.com",
			wantErr:      true,
		},
		{
			name:         "empty target",
			key:          1,
			shortcutName: "youtube",
			shortcutType: "url",
			target:       "   ",
			wantErr:      true,
		},
		{
			name:         "invalid URL",
			key:          1,
			shortcutName: "youtube",
			shortcutType: "url",
			target:       "banana",
			wantErr:      true,
		},
		{
			name:         "nonexistent folder",
			key:          1,
			shortcutName: "folder",
			shortcutType: "folder",
			target:       filepath.Join(os.TempDir(), "does-not-exist"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateShortcut(
				tt.key,
				tt.shortcutName,
				tt.shortcutType,
				tt.target,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuildIndex(t *testing.T) {
	shortcuts := []shortcut{
		{
			ID:     "abc123",
			Key:    1,
			Name:   "youtube",
			Type:   "url",
			Target: "https://youtube.com",
		},
		{
			ID:     "xyz789",
			Key:    2,
			Name:   "chatgpt",
			Type:   "url",
			Target: "https://chatgpt.com",
		},
	}

	index := buildIndex(shortcuts)

	entry, ok := index["abc123"]
	if !ok {
		t.Fatal("expected abc123 to exist")
	}

	if entry.Name != "youtube" {
		t.Errorf("expected youtube, got %s", entry.Name)
	}

	entry, ok = index["2"]
	if !ok {
		t.Fatal("expected key 2 to exist")
	}

	if entry.Name != "chatgpt" {
		t.Errorf("expected chatgpt, got %s", entry.Name)
	}

	_, ok = index["does-not-exist"]
	if ok {
		t.Error("expected nonexistent ID to not exist")
	}
}

func TestAdd(t *testing.T) {
	dir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	err = os.WriteFile("shortcuts.json", []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Add(1, "youtube", "url", "https://youtube.com")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	file, err := os.Open("shortcuts.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var shortcuts []shortcut
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		t.Fatal(err)
	}

	if len(shortcuts) != 1 {
		t.Fatalf("expected 1 shortcut, got %d", len(shortcuts))
	}

	if shortcuts[0].ID == "" {
		t.Error("expected generated ID")
	}

	if shortcuts[0].Key != 1 {
		t.Errorf("expected key 1, got %d", shortcuts[0].Key)
	}

	if shortcuts[0].Name != "youtube" {
		t.Errorf("expected youtube, got %s", shortcuts[0].Name)
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	data := `[
		{"ID":"abc","Key":1,"Name":"youtube","Type":"url","Target":"https://youtube.com"},
		{"ID":"xyz","Key":2,"Name":"chatgpt","Type":"url","Target":"https://chatgpt.com"}
	]`

	err = os.WriteFile("shortcuts.json", []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Remove(1)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	file, err := os.Open("shortcuts.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var shortcuts []shortcut
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		t.Fatal(err)
	}

	if len(shortcuts) != 1 {
		t.Fatalf("expected 1 shortcut, got %d", len(shortcuts))
	}

	if shortcuts[0].Key != 2 {
		t.Errorf("expected remaining shortcut to have key 2, got %d", shortcuts[0].Key)
	}
}

func TestEdit(t *testing.T) {
	dir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	data := `[
		{"ID":"abc123","Key":1,"Name":"youtube","Type":"url","Target":"https://youtube.com"}
	]`

	err = os.WriteFile("shortcuts.json", []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Edit(1, "name", "my youtube")
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	file, err := os.Open("shortcuts.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var shortcuts []shortcut
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		t.Fatal(err)
	}

	if shortcuts[0].Name != "my youtube" {
		t.Errorf("expected edited name, got %s", shortcuts[0].Name)
	}

	if shortcuts[0].ID != "abc123" {
		t.Errorf("expected ID to remain abc123, got %s", shortcuts[0].ID)
	}
}

func TestIDGeneration(t *testing.T) {
	dir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	err = os.WriteFile("shortcuts.json", []byte("[]"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = Add(1, "youtube", "url", "https://youtube.com")
	if err != nil {
		t.Fatal(err)
	}

	err = Add(2, "chatgpt", "url", "https://chatgpt.com")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open("shortcuts.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var shortcuts []shortcut
	err = json.NewDecoder(file).Decode(&shortcuts)
	if err != nil {
		t.Fatal(err)
	}

	if shortcuts[0].ID == "" || shortcuts[1].ID == "" {
		t.Error("expected both shortcuts to have IDs")
	}

	if shortcuts[0].ID == shortcuts[1].ID {
		t.Error("expected unique IDs")
	}
}
