package models

import (
	"strings"
)

type Cipher struct {
	Name        string
	Length      int
	Cols        int
	Rows        int
	Symbols     []rune // The actual cipher
	Translation string
	Solved      bool
	Key         map[rune]rune
	LettersKey  map[string][]string
	SymbolsKey  map[string][]string
}

type Character struct {
	Symbol    string
	Letter    string
	EndOfWord bool
}

func NewCipher(str string, cols int, solved bool, key map[rune]rune) *Cipher {
	c := new(Cipher)
	c.Symbols = []rune(str)
	c.Length = len(str)
	c.Key = key
	c.doTranslation()

	c.Cols = cols
	c.Rows = c.Length / c.Cols
	c.Solved = solved // boolean

	return c
}

func (c *Cipher) doTranslation() {

	temp := make([]rune, c.Length)

	for s, symbol := range c.Symbols {
		temp[s] = c.Key[symbol]
	}

	c.Translation = string(temp)
}

func (c *Cipher) GetSymbol(index int) string {
	return string(c.Symbols[index])
}

func (c *Cipher) GetRows() int {
	return c.Rows
}

func (c *Cipher) GetCols() int {
	return c.Cols
}

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

// Display Letters Key
func (c *Cipher) DisplayKey() map[string]string {

	stringMap := make(map[string]string, c.Length)

	for symbol, letter := range c.Key {
		stringMap[string(symbol)] = string(letter)
	}

	return stringMap
}

// Build Symbols Key
func (c *Cipher) BuildLettersKey() {

	t := []byte(strings.Replace(c.Translation, " ", "", -1))

	if c.Key == nil {
		c.LettersKey = make(map[string][]string)
	}

Loop:
	for i, letter := range t {
		nS := string(c.Symbols[i])
		for _, existing := range c.LettersKey[string(letter)] {
			if existing == nS {
				continue Loop
			}
		}
		c.LettersKey[string(letter)] = append(c.LettersKey[string(letter)], nS)
	}
}
