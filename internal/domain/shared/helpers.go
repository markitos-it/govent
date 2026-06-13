package shared

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

func RandomString(howManyLetters ...int) string {
	length := 10
	if len(howManyLetters) > 0 {
		length = howManyLetters[0]
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(), RandomString())
}

func IsUUIDv4(uuid string) bool {
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(uuidRegex, uuid)

	return err == nil && matched
}

func UUIDv4() string {
	var uuid [16]byte
	rand.Read(uuid[:])

	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func RandomWord(howManyChars ...int) string {
	length := 10
	if len(howManyChars) > 0 {
		length = howManyChars[0]
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

func RandomSlug() string {
	const firstChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const middleChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-"
	const lastChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	num, _ := rand.Int(rand.Reader, big.NewInt(20))
	length := int(num.Int64()) + 2

	result := make([]byte, length)

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(firstChars))))
	result[0] = firstChars[n.Int64()]

	for i := 1; i < length-1; i++ {
		n, _ = rand.Int(rand.Reader, big.NewInt(int64(len(middleChars))))
		result[i] = middleChars[n.Int64()]
	}

	n, _ = rand.Int(rand.Reader, big.NewInt(int64(len(lastChars))))
	result[length-1] = lastChars[n.Int64()]

	return string(result)
}

func RandomPersonalName() string {
	minWords, maxWords := 1, 6
	wordCount, _ := rand.Int(rand.Reader, big.NewInt(int64(maxWords-minWords+1)))
	wordCount = big.NewInt(wordCount.Int64() + int64(minWords))
	minLength, maxLength := 3, 150

	var result []string
	currentLength := 0

	for currentLength < minLength || currentLength+int(wordCount.Int64())-1 > maxLength {
		result = nil
		currentLength = 0
		for i := 0; i < int(wordCount.Int64()); i++ {
			num, _ := rand.Int(rand.Reader, big.NewInt(8))
			word := RandomString(int(num.Int64()) + 3)
			word = strings.ToLower(word)
			result = append(result, word)
			currentLength += len(word)
		}
	}

	return strings.Join(result, " ")
}

func Slug() string {
	minLength, maxLength := 3, 100

	var result []string
	currentLength := 0

	for currentLength < minLength || currentLength+len(result)-1 > maxLength {
		result = nil
		num, _ := rand.Int(rand.Reader, big.NewInt(6))
		wordCount := int(num.Int64()) + 1
		for i := 0; i < wordCount; i++ {
			num, _ := rand.Int(rand.Reader, big.NewInt(8))
			word := RandomString(int(num.Int64()) + 3)
			word = strings.ToLower(word)
			result = append(result, word)
			currentLength += len(word)
		}
	}

	return strings.Join(result, "-")
}
