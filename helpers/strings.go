package helpers

import "math/rand"

func GetReversedString(s string) string {
	result := ""
	for i := len(s) - 1; i >= 0; i-- {
		result += string(s[i])
	}
	return result
}

func GetRandomString(letters []rune, length uint) string {
	res := make([]rune, length)
	for i := range res {
		res[i] += letters[rand.Intn(len(letters))]
	}
	return string(res)
}
