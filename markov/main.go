package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	MIN_SENTENCE_LEN int = 5
)

func main() {
	inPath := flag.String("i", "", "Document that will source the markov chain")
	outPath := flag.String("o", "", "The output file")
	flag.Parse()

	t := time.Now()
	inFile, err := os.Open(*inPath)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist")
		os.Exit(1)
	} else if err != nil {
		fmt.Println("Could not open file: ", err)
		os.Exit(2)
	}
	defer inFile.Close()

	inContent, err := io.ReadAll(inFile)
	if err != nil {
		fmt.Println("Could not read file: ", err)
		os.Exit(3)
	}

	data := sanitize(string(inContent))
	mapChain(data)

	outContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Could not marshal json: ", err)
		os.Exit(4)
	}

	outFile, err := os.Create(*outPath)
	if err != nil {
		fmt.Println("Could not create file: ", err)
		os.Exit(5)
	}

	defer outFile.Close()

	_, err = outFile.Write(outContent)
	if err != nil {
		fmt.Println("Could not write output: ", err)
		os.Exit(6)
	}

	fmt.Printf("Built chain in %s\n", time.Since(t).String())
}

func sanitize(input string) map[string][]string {
	input = strings.ToLower(input)
	splitRe := regexp.MustCompile("[\n.!?]+")
	textRe := regexp.MustCompile("[^a-zA-Z']+")
	apostrRe := regexp.MustCompile("'")
	spaceRe := regexp.MustCompile(" +")

	output := make(map[string][]string)

	sentences := splitRe.Split(input, -1)
	space := []byte(" ")
	for _, sentence := range sentences {
		cleaned := string(spaceRe.ReplaceAll(textRe.ReplaceAll([]byte(sentence), space), space))
		cleaned = apostrRe.ReplaceAllString(cleaned, "")

		words := strings.Split(cleaned, " ")
		if MIN_SENTENCE_LEN > len(words) {
			continue
		}
		last := ""
		for _, word := range words {
			if last == "" {
				last = word
				continue
			}
			if !slices.Contains(output[last], word) {
				output[last] = append(output[last], word)
			}
			last = word
		}
	}
	return output
}

func mapChain(data map[string][]string) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	cleanDeadLinks(keys, data)
	cleanEmptyKeys(keys, data)

	slices.Reverse(keys)

	cleanDeadLinks(keys, data)
	cleanEmptyKeys(keys, data)
}

func cleanDeadLinks(keys []string, data map[string][]string) {
	for _, key := range keys {
		value, ok := data[key]
		if !ok {
			continue
		}
		words := make([]string, 0, len(value))
		for _, word := range value {
			word = strings.TrimSpace(word)
			if _words, ok := data[word]; ok { // Skips nested dead links
				if len(_words) != 0 {
					words = append(words, word)
				}
			}
		}
		data[key] = words
	}

}

func cleanEmptyKeys(keys []string, data map[string][]string) {
	for _, key := range keys {
		value, ok := data[key]
		if !ok {
			continue
		}
		if len(value) == 0 {
			delete(data, key)
		}
	}
}
