package words

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"rand"
	"time"
)

type WordSet string

const (
	EnPop200 WordSet = "English200Popular"
)

var availableWordSets = []WordSet{enPop200}
var wordSetToFilename = map[WordSet]string{enPop200: "en-pop-200.json"}

type WordSetData struct {
	name  string
	size  int
	words []string
}

func ShuffleWordSet(set WordSet) ([]string, error) {
	var setFilename string
	if Contains[WordSet](availableWordSets, set) {
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

	var words []string = payload.words
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	return words, nil
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
