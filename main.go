package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/VincentBrodin/tcli/app"
)

const WORD_FILE_URL string = "https://raw.githubusercontent.com/VincentBrodin/godo/refs/heads/main/words.json"

func main() {
	length := flag.Int("l", 10, "How length of the test")
	flag.Parse()

	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Could not find executable path:", err)
		os.Exit(1)
	}
	execDir := filepath.Dir(execPath)
	wordFile := filepath.Join(execDir, "words.json")

	// Load the word chain
	file, err := os.Open(wordFile)
	if os.ErrNotExist == err {
		fmt.Println("Could not find word file, grabbing it from ", WORD_FILE_URL)
		if err := getWordFile(wordFile); err != nil {
			fmt.Println("Failed to get word file: ", err)
			os.Exit(2)
		} else {
			file, err = os.Open(wordFile)
		}
	}
	if err != nil {
		fmt.Println("Could not open words.json: ", err)
		os.Exit(3)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Could read words.json: ", err)
		os.Exit(4)
	}

	words := make(map[string][]string)

	if err := json.Unmarshal(data, &words); err != nil {
		fmt.Println("Could unmarshal words.json: ", err)
		os.Exit(5)
	}

	a := app.App{
		Words:  words,
		Length: *length,
	}

	if err := a.Run(); err != nil {
		os.Exit(6)
	}
}

func getWordFile(path string) error {
	resp, err := http.Get(WORD_FILE_URL)
	if err != nil {
		fmt.Println("Failed to grab word file: ", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create words.json: ", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Failed to write words.json: ", err)
	}

	return nil
}
