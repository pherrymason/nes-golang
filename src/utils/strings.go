package utils

import (
	"encoding/hex"
	"regexp"
)

func StringSplitByRegex(content string) []string {
	re := regexp.MustCompile("[\\s]{2,}")

	return re.Split(content, -1)
}

func HexStringToByteArray(field string) []byte {
	result, err := hex.DecodeString(field)
	if err != nil {
		panic(err)
	}
	return result
}
