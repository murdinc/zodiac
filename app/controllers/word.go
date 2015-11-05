package controllers

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type WordList []Word

type Word struct {
	Value string
	Count int
}

func (w WordList) Len() int      { return len(w) }
func (w WordList) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

type ByCount struct{ WordList }

func (c ByCount) Less(i, j int) bool { return c.WordList[i].Count > c.WordList[j].Count }

func GetWordList(wordListDir string) (*WordList, error) {

	//wordListDir := revel.AppPath + "/words/"

	files, err := ioutil.ReadDir(wordListDir)
	if err != nil {
		return &WordList{}, err
	}

	wordList := WordList{}

	for _, file := range files {

		if strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {

			wordFile, err := Open(wordListDir + file.Name())
			if err != nil {
				return &WordList{}, err
			}

			wordList = append(wordList, wordFile...)

		}
		// Deduper 4000™
		wordList.removeDuplicates()
	}

	return &wordList, nil

}

func (w *WordList) removeDuplicates() {
	found := make(map[string]bool)
	j := 0
	for i, x := range *w {
		if !found[strings.ToLower(x.Value)] {
			found[strings.ToLower(x.Value)] = true
			(*w)[j] = (*w)[i]
			j++
		}
	}
	*w = (*w)[:j]
}

func Open(fileName string) (WordList, error) {

	wordList := WordList{}

	f, err := os.Open(fileName)

	if err != nil {
		return WordList{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) > 3 {
			wordList = append(wordList, Word{Value: scanner.Text()})
		}
	}

	if err := scanner.Err(); err != nil {
		return WordList{}, err
	}

	return wordList, nil
}
