package words

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type WordSet string

const (
	EnPop200          WordSet = "English200Popular"
	wordsPerLineLimit         = 10
)

var availableWordSets = []WordSet{EnPop200}
var wordSetToFilename = map[WordSet]string{EnPop200: "word_sets/en-pop-200.json"}

type WordSetData struct {
	Name  string   `json:"name"`
	Size  int      `json:"size"`
	Words []string `json:"words"`
}

func ShuffleWordSet(set WordSet) ([]rune, error) {
	var setFilename string
	if Contains(availableWordSets, set) {
		setFilename = wordSetToFilename[set]
	}
	if setFilename == "" {
		log.Fatal("Unknown wordset")
		return nil, errors.New("Unknown wordset")
	}
	content, err := ioutil.ReadFile(setFilename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return nil, err
	}

	var payload WordSetData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return nil, err
	}

	var words []string = payload.Words
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	var runes []rune
	var wordsInLine = 0
	for _, st := range words {
		chars := []rune(st)
		chars = append(chars, ' ')
		if wordsInLine >= wordsPerLineLimit {
			chars = append(chars, '\n')
			wordsInLine = 0
		}
		runes = append(runes, chars...)
		wordsInLine++
	}
	return runes, nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
