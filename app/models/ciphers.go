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
}

type Character struct {
	Symbol    string
	Letter    string
	EndOfWord bool
}

func NewCipher(str string, cols int, solved bool, translation string) *Cipher {
	c := new(Cipher)
	c.Symbols = []rune(str)
	c.Length = len(str)
	c.Cols = cols
	c.Rows = c.Length / c.Cols
	c.Solved = solved
	c.Translation = translation
	return c
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

func (c *Cipher) Display() [][]Character {

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
