package controllers

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/revel/revel"
)

type Cipher struct {
	Name            string
	Length          int
	Cols            int
	Rows            int
	Symbols         []rune // The actual cipher
	Translation     string // plain text translation
	Solved          bool
	Key             map[rune]rune
	SymbolCount     map[rune]int
	FoundWordsTotal int
	FoundWords      map[string]Word // map[word]count
	WordLengths     map[int]int     // map[wordLength]count
}

type Character struct {
	Symbol    string
	Letter    string
	EndOfWord bool
}

type Letter struct {
	letter    rune
	reuseable bool
}

// Creates a new Cipher struct and returns it.
func NewCipher(str string) *Cipher {
	c := new(Cipher)

	c.Symbols = []rune(str)
	c.Length = len(str)
	c.SymbolCount = make(map[rune]int)

	return c
}

// Set the Bool flag in the Cipher struct
func (c *Cipher) SetSolved(solved bool) {
	c.Solved = solved
}

// Set the Key map in the Cipher struct
func (c *Cipher) SetKey(key map[rune]rune) {
	c.Key = key
	c.doTranslation()

	if len(c.SymbolCount) == 0 {
		for _, letter := range c.Key {
			c.incrementSymbolCount(letter)
		}
	}
}

// Do the translation from Symbol to Letter
func (c *Cipher) doTranslation() {

	temp := make([]rune, c.Length)

	for s, symbol := range c.Symbols {
		temp[s] = c.Key[symbol]
	}

	c.Translation = string(temp)
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

// Find a word in a cipher
func (c *Cipher) FindWord() int {
	re := regexp.MustCompile("DO")
	matches := re.FindAllString(c.Translation, -1)
	return len(matches)
}

// Displays the Cipher
func (c *Cipher) DisplayCipher() [][]Character {

	t := []byte(strings.Replace(c.Translation, " ", "", -1))

	d := make([][]Character, c.Rows) // Outer Lasyer: Rows
	for i := range d {
		d[i] = make([]Character, c.Cols) // Inner Layer; Cols
		for j := range d[i] {
			d[i][j].Symbol = string(c.Symbols[j+((i)*c.Cols)])
			d[i][j].Letter = string(t[j+((i)*c.Cols)])
		}
	}

	return d
}

// Build Symbols Key
func (c *Cipher) DisplaySymbolsKey() map[string]string {

	symbolsKey := make(map[string]string, c.Length)

	for symbol, letter := range c.Key {
		symbolsKey[string(symbol)] = string(letter)
	}

	return symbolsKey
}

// Build Letters Key
func (c *Cipher) DisplayLettersKey() map[string][]string {

	t := []byte(strings.Replace(c.Translation, " ", "", -1))

	lettersKey := make(map[string][]string)

Loop:
	for i, letter := range t {
		nS := string(c.Symbols[i])
		for _, existing := range lettersKey[string(letter)] {
			if existing == nS {
				continue Loop
			}
		}
		lettersKey[string(letter)] = append(lettersKey[string(letter)], nS)
	}

	return lettersKey
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
func (c *Cipher) RandomKey(cipherString []rune) map[rune]rune {

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
		return nil
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

		c.incrementSymbolCount(chosenOne.letter)

		if !chosenOne.reuseable || c.SymbolCount[chosenOne.letter] >= revel.Config.IntDefault("cipher.maxSymbols", 4) {
			letters = append(letters[:0], letters[1:]...) // stop using this letter
		}

	}

	return randomkey
}

func (c *Cipher) incrementSymbolCount(letter rune) {

	switch c.SymbolCount[letter] {
	case 0:
		c.SymbolCount[letter] = 1
	default:
		c.SymbolCount[letter]++
	}

}

func (c *Cipher) GetSymbolCount(letter string) int {
	return c.SymbolCount[rune(letter[0])]
}

// Shuffle function used in random key generation
func Shuffle(a []Letter) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Returns the basic (my version) Key for the Z408 Cipher
func (c *Cipher) Z408Key() map[rune]rune {
	return map[rune]rune{
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
}
