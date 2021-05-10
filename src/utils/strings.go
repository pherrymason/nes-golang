package utils

import (
	"encoding/hex"
	"regexp"
	"strconv"
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

func ByteToHex(value byte) string {
	buf := []byte{'0', '0', '0', '0', 4 + 16: 0}
	buf = strconv.AppendInt(buf[:4], int64(value), 16)

	return string(buf[len(buf)-4:])
}
