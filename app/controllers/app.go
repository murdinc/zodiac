package controllers

import (
	"zodiac/app/models"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	cipher := models.NewCipher(
		revel.Config.StringDefault("cipher.340.raw", ""),
		revel.Config.IntDefault("cipher.340.cols", 10),
		revel.Config.BoolDefault("cipher.340.solved", false),
		map[rune]rune{},
	)
	return c.Render(cipher)
}

func (c App) Z408() revel.Result {

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

	cipher := models.NewCipher(
		revel.Config.StringDefault("cipher.408.raw", ""),
		revel.Config.IntDefault("cipher.408.cols", 10),
		revel.Config.BoolDefault("cipher.408.solved", false),
		key,
	)
	return c.Render(cipher)
}

func (c App) Z340() revel.Result {

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
	cipher := models.NewCipher(
		revel.Config.StringDefault("cipher.340.raw", ""),
		revel.Config.IntDefault("cipher.340.cols", 10),
		revel.Config.BoolDefault("cipher.340.solved", false),
		key,
	)
	return c.Render(cipher)
}
