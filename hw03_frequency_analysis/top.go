package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	wordsCount := map[string]int{}
	words := make([]string, 0, len(wordsCount))

	// Заполняем словарь (слово - кол-во в тексте)
	for _, v := range strings.Fields(s) {
		key := prepareKey(v)
		wordsCount[key]++
	}

	// Преобразуем словарь в слайс со словами
	for word := range wordsCount {
		if word == "-" {
			continue
		}
		words = append(words, word)
	}

	// Сортируем слайс по кол-ву слов в тексте
	sort.Slice(words, func(i int, j int) bool {
		// Если кол-во слов равно, то сравниваем сами слова (лексикографическая сортировка)
		if wordsCount[words[i]] == wordsCount[words[j]] {
			return words[i] < words[j]
		}
		return wordsCount[words[i]] > wordsCount[words[j]]
	})

	// На случай, если слов в тексте меньше 10, высчитываем длину итогового слайса
	sliceLen := len(words)
	if sliceLen > 10 {
		sliceLen = 10
	}

	return words[:sliceLen]
}

func prepareKey(s string) string {
	return strings.ToLower(strings.Trim(s, ",!.'"))
}
