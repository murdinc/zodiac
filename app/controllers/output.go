package controllers

import (
	"strings"

	"github.com/revel/revel"
)

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

func (c *Cipher) KeyCount() int {
	count, err := GetKeyCount(c)
	if err != nil {
		revel.ERROR.Print(err)
	}

	return count

}

// Build Symbols Key
func (c *Cipher) DisplaySymbolsKey() map[string]string {

	symbolsKey := make(map[string]string, c.Length)

	for _, character := range c.Key {
		symbolsKey[string(character.Symbol)] = string(character.Letter)
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
