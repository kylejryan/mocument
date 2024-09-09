package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateRandomData(field FieldSchema) interface{} {
	switch field.Type {
	case "string":
		return generateRandomString(field)
	case "float":
		return generateRandomFloat(field)
	case "integer":
		return generateRandomInteger(field)
	default:
		return nil
	}
}

// Updated to return interface{} so it can be assigned in the map.
func generateRandomString(field FieldSchema) interface{} {
	if len(field.Options) > 0 {
		return field.Options[rand.Intn(len(field.Options))]
	}
	if field.Format != "" {
		return generateFormattedString(field.Format)
	}
	return randomString(10) // Default random string
}

func generateRandomFloat(field FieldSchema) interface{} {
	return field.Min + rand.Float64()*(field.Max-field.Min)
}

func generateRandomInteger(field FieldSchema) interface{} {
	return int(field.Min) + rand.Intn(int(field.Max-field.Min))
}

func generateFormattedString(format string) string {
	result := ""
	for _, c := range format {
		if c == '#' {
			result += fmt.Sprintf("%d", rand.Intn(10))
		} else {
			result += string(c)
		}
	}
	return result
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
