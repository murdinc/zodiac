package models

type Solution struct {
	ID     int           // hex string?
	Key    map[rune]rune // The Key for this particular solution
	String string
}

func (s *Solution) GetKey() map[rune]rune {
	return s.Key
}
