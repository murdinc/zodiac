package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Cipher struct {
	Name            string
	Length          int
	Cols            int
	Rows            int
	Symbols         []rune // The actual cipher
	Translation     string // plain text translation
	Solved          bool
	Key             []Character
	KeyID           string
	SymbolCount     map[string]int
	FoundWordsTotal int
	FoundWords      WordList
	WordLengths     map[int]int // map[wordLength]count
	KillCount       int
	Source          string
}

type KeyDoc struct {
	CipherName      string
	KeyID           string
	Translation     string
	Timestamp       int64
	KillCount       int
	Key             []Character
	FoundWordsTotal int
	FoundWords      WordList
	Source          string
}

type Character struct {
	Symbol string
	Letter string
}

type Key []Character

func (k Key) Len() int      { return len(k) }
func (k Key) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

type BySymbol struct{ Key }

func (s BySymbol) Less(i, j int) bool { return s.Key[i].Symbol < s.Key[j].Symbol }

type Letter struct {
	letter    rune
	reuseable bool
}

// Creates a new Cipher struct and returns it.
func NewCipher(name string, str string) *Cipher {
	c := new(Cipher)

	c.Name = name
	c.Symbols = []rune(str)
	c.Length = len(str)
	c.SymbolCount = make(map[string]int)

	return c
}

// Generate MD5 hash from key data
func HashKey(key []Character) string {

	sort.Sort(BySymbol{key})

	// Build a string of the sorted keys
	sortedStringKey := ""
	for _, character := range key {
		sortedStringKey = sortedStringKey + character.Symbol + character.Letter
	}

	hasher := md5.New()
	hasher.Write([]byte(sortedStringKey))

	id := hex.EncodeToString(hasher.Sum(nil))

	return id
}

// Set the Solved flag in the Cipher struct
func (c *Cipher) SetSolved(solved bool) {
	c.Solved = solved
}

// Set the Key in the Cipher struct from a given key map
func (c *Cipher) setKeyFromKeyMap(keyMap map[rune]rune) {

	key := make([]Character, 0)

	for symbol, letter := range keyMap {
		key = append(key, Character{Letter: string(letter), Symbol: string(symbol)})
	}

	// Count Symbol Occurance
	if len(c.SymbolCount) == 0 {
		for _, character := range c.Key {
			c.incrementSymbolCount(character.Letter)
		}
	}

	c.Key = key
	c.KeyID = HashKey(key)
	c.doTranslation()
}

// Set the Key in the Cipher struct from a given key map
func (c *Cipher) SetKeyFromKeyDoc(keyDoc KeyDoc) {
	c.Key = keyDoc.Key
	c.Translation = keyDoc.Translation
	c.KeyID = keyDoc.KeyID
	c.FoundWordsTotal = keyDoc.FoundWordsTotal
	c.FoundWords = keyDoc.FoundWords
	c.KillCount = keyDoc.KillCount
	c.Source = keyDoc.Source

	// Count Symbol Occurance
	if len(c.SymbolCount) == 0 {
		for _, character := range c.Key {
			c.incrementSymbolCount(character.Letter)
		}
	}
}

// Do the translation from Symbol to Letter
func (c *Cipher) doTranslation() {

	temp := ""

	for _, symbol := range c.Symbols {
	Loop:
		for _, character := range c.Key {
			if character.Symbol == string(symbol) {
				temp = temp + character.Letter
				break Loop
			}
		}

	}

	c.Translation = temp
	c.KillCount = strings.Count(c.Translation, "KILL")

}

// Gets a Symbol given the index of it in the string
func (c *Cipher) GetSymbol(index int) string {
	return string(c.Symbols[index])
}

// Set the number of Cols (and Rows)
func (c *Cipher) SetCols(cols int) {
	c.Cols = cols
	c.Rows = c.Length / c.Cols
}

// Get the number of Cols
func (c *Cipher) GetCols() int {
	return c.Cols
}

// Get the number of Rows
func (c *Cipher) GetRows() int {
	return c.Rows
}

// Scan for words in our word list
func (c *Cipher) ScanForWords(words *WordList) {

	foundWordsTotal := 0
	foundWords := WordList{}

	for _, word := range *words {
		count := strings.Count(c.Translation, strings.ToUpper(word.Value))
		if count > 0 {
			foundWordsTotal += count
			foundWords = append(foundWords, Word{Value: word.Value, Count: count})
		}
	}

	sort.Sort(ByCount{foundWords})
	c.FoundWords = foundWords
	c.FoundWordsTotal = foundWordsTotal

}

func BuildLetters() []Letter {
	letters := make([]Letter, 26)
	alphabet := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	reuseableLetters := []rune{'E', 'T', 'A', 'O', 'N', 'I', 'R', 'S', 'L', 'D', 'F', 'H'} // maybe expand more?

Loop:
	for i, letter := range alphabet {
		for _, rL := range reuseableLetters {
			if rL == letter {
				letters[i] = Letter{letter: letter, reuseable: true}
				continue Loop
			}
		}
		letters[i] = Letter{letter: letter, reuseable: false}
	}

	return letters

}

// Generates a Random Key, given a rune array of the cipher
func (c *Cipher) RandomKey(cipherString []rune, maxSymbols int) (KeyDoc, error) {

	letters := BuildLetters()

	occur := [256]byte{}
	order := make([]rune, 0, 256)

	n := 0

	for i := 0; i < len(cipherString); i++ {
		b := cipherString[i]
		switch occur[b] {
		case 0:
			occur[b] = 1
			order = append(order, b)
			n++
		case 1:
			occur[b]++
			n--
		}
	}

	if n == 0 {
		return KeyDoc{}, errors.New("ERROR GENERATING RANDOM KEY!! ")
	}

	symbols := make([]rune, 0, n)

	for _, b := range order {
		symbols = append(symbols, rune(b))
	}

	randomkey := make(map[rune]rune)

	for _, symbol := range symbols {

		rand.Seed(time.Now().UnixNano())
		Shuffle(letters)

		chosenOne := letters[0]

		randomkey[symbol] = chosenOne.letter

		c.incrementSymbolCount(string(chosenOne.letter))

		if !chosenOne.reuseable || c.SymbolCount[string(chosenOne.letter)] >= maxSymbols {
			letters = append(letters[:0], letters[1:]...) // stop using this letter
		}

	}

	c.setKeyFromKeyMap(randomkey)

	keyDoc := KeyDoc{
		Translation:     c.Translation,
		Key:             c.Key,
		KeyID:           c.KeyID,
		CipherName:      c.Name,
		Timestamp:       time.Now().Unix(),
		FoundWordsTotal: c.FoundWordsTotal,
		FoundWords:      c.FoundWords,
		KillCount:       c.KillCount,
		Source:          "zodiac-server",
	}

	return keyDoc, nil
}

func (c *Cipher) incrementSymbolCount(letter string) {

	switch c.SymbolCount[letter] {
	case 0:
		c.SymbolCount[letter] = 1
	default:
		c.SymbolCount[letter]++
	}

}

func (c *Cipher) GetSymbolCount(letter string) int {
	return c.SymbolCount[string(letter[0])]
}

// Shuffle function used in random key generation
func Shuffle(a []Letter) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Returns the basic (my version) Key for the Z408 Cipher
func (c *Cipher) Z408Solution() (KeyDoc, error) {
	key := map[rune]rune{
		'!':  'A',
		'#':  'A',
		'$':  'L',
		'%':  'L',
		'&':  'Y',
		'*':  'S',
		'+':  'E',
		'/':  'K',
		'0':  'T',
		'5':  'E',
		'6':  'H',
		'7':  'N',
		'9':  'D',
		'@':  'I',
		'A':  'W',
		'B':  'L',
		'D':  'N',
		'E':  'E',
		'F':  'S',
		'G':  'A',
		'H':  'T',
		'J':  'F',
		'K':  'S',
		'L':  'T',
		'M':  'H',
		'N':  'E',
		'O':  'N',
		'P':  'I',
		'Q':  'F',
		'R':  'G',
		'S':  'A',
		'T':  'O',
		'U':  'I',
		'V':  'B',
		'W':  'E',
		'X':  'O',
		'Y':  'U',
		'Z':  'E',
		'[':  'T',
		'\\': 'R',
		']':  'O',
		'c':  'V',
		'd':  'O',
		'e':  'C',
		'f':  'D',
		'j':  'X',
		'k':  'I',
		'l':  'A',
		'p':  'E',
		'q':  'M',
		'r':  'R',
		't':  'R',
		'v':  'N',
		'~':  'P',
	}

	c.setKeyFromKeyMap(key)

	// Count Symbol Occurance
	if len(c.SymbolCount) == 0 {
		for _, character := range c.Key {
			c.incrementSymbolCount(character.Letter)
		}
	}

	keyDoc := KeyDoc{
		Translation:     c.Translation,
		Key:             c.Key,
		KeyID:           c.KeyID,
		CipherName:      c.Name,
		Timestamp:       time.Now().Unix(),
		FoundWordsTotal: c.FoundWordsTotal,
		FoundWords:      c.FoundWords,
		KillCount:       c.KillCount,
		Source:          "zodiac-server",
	}

	return keyDoc, nil

}
