package models

type Cipher struct {
	Name     string
	Length   int
	Cols     int
	Rows     int
	Symbols  []rune // The actual cipher
	Solution Solution
}

func NewCipher(str string, cols int) *Cipher {
	c := new(Cipher)
	c.Symbols = []rune(str)
	c.Length = len(str)
	c.Cols = cols
	c.Rows = c.Length / c.Cols
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

func (c *Cipher) Display() [][]string {

	d := make([][]string, c.Rows) // Outer Lasyer: Rows
	for i := range d {
		d[i] = make([]string, c.Cols) // Inner Layer; Cols
		for j := range d[i] {
			d[i][j] = string(c.Symbols[j+((i)*c.Cols)])
		}
	}

	return d
}
