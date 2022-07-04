package helpers

import (
	"fmt"
	"math/rand"
)

func ParsePercent(p float64) string {
	return fmt.Sprintf("%.0f", (p*100)) + "%"
}

func ParseKey(k int) string {
	switch k {
	case 0:
		return "C"
	case 1:
		return "C♯"
	case 2:
		return "D"
	case 3:
		return "D♯"
	case 4:
		return "E"
	case 5:
		return "F"
	case 6:
		return "F♯"
	case 7:
		return "G"
	case 8:
		return "G♯"
	case 9:
		return "A"
	case 10:
		return "A♯"
	case 11:
		return "B"
	default:
		return "??"
	}
}

func ParseTempo(t float64) string {
	return fmt.Sprintf("%.0f BPM", t)
}

func ParseLoudness(l float64) string {
	return fmt.Sprintf("%.0f db", l)
}

func GenerateRandomString(n int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
