package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordInfo struct {
	word  string
	count int
}

var sep = regexp.MustCompile(`[:.,!?;"'\s]+`)

func Top10(text string) []string {
	return Top(text, 10)
}

func Top(text string, n int) []string {
	if n < 0 {
		return nil
	}

	countMap := getCountMap(text)
	return getSortedWords(countMap, n)
}

func getCountMap(text string) map[string]int {
	countMap := map[string]int{}

	for _, word := range sep.Split(text, -1) {
		if word == "" || word == "-" || word == "–" {
			continue
		}

		countMap[strings.ToLower(word)]++
	}

	return countMap
}

func getSortedWords(countMap map[string]int, n int) []string {
	slice := make([]wordInfo, len(countMap))
	i := 0
	for word, count := range countMap {
		slice[i] = wordInfo{word, count}
		i++
	}

	sort.Slice(slice, func(i int, j int) bool {
		if slice[i].count == slice[j].count {
			return slice[i].word < slice[j].word
		}

		return slice[i].count > slice[j].count
	})

	var quantity int
	if n > len(countMap) {
		quantity = len(countMap)
	} else {
		quantity = n
	}

	sortedWords := make([]string, quantity)
	for i := 0; i < quantity; i++ {
		sortedWords[i] = slice[i].word
	}

	return sortedWords
}
