package main

import (
	"io/ioutil"
	"regexp"
	"fmt"
	"strings"
)

func train(training_data string) map[string]int {
	NWORDS := make(map[string]int)
	pattern := regexp.MustCompile("[a-z]+")
	if content, err := ioutil.ReadFile(training_data); err == nil {
 		for _, w := range pattern.FindAllString(strings.ToLower(string(content)), -1) {
			NWORDS[w]++;
		}
	} else {
		panic("Failed loading training data.  Get it from http://norvig.com/big.txt.")
	}
	return NWORDS
}

func edits1(word string, ch chan string) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	type Pair struct{a, b string}
	var splits []Pair
	for i := 0; i < len(word) + 1; i++ {
		splits = append(splits, Pair{word[:i], word[i:]}) }

	for _, s := range splits { // deletes
		if len(s.b) > 0 { ch <- s.a + s.b[1:] }
		if len(s.b) > 1 { ch <- s.a + string(s.b[1]) + string(s.b[0]) + s.b[2:] }
		for _, c := range alphabet { if len(s.b) > 0 { ch <- s.a + string(c) + s.b[1:] }}
		for _, c := range alphabet { ch <- s.a + string(c) + s.b }
	}
}

func edits2(word string, ch chan string) {
	edits1ch := make(chan string, 10)
	go func() {
		edits1(word, edits1ch)
		ch <- ""}()
	for e1 := range edits1ch {
		if e1 == "" { break }
		edits1(e1, ch)
	}
}

func correct(word string, NWORDS map[string]int) string {
	ch := make(chan string)
	go func() {
		ch <- word
		edits1(word, ch)
		edits2(word, ch)
		ch <- ""
	}()
	maxFreq := 0
	correction := ""
	for word := range ch {
		if word == "" { return correction }
		if freq, present := NWORDS[word]; present && freq > maxFreq {
			maxFreq, correction = freq, word
		}
	}
	return ""
}

func main() {
	model := train("big.txt")
	fmt.Println(correct("speling", model))
	fmt.Println(correct("korrecter", model))
}
