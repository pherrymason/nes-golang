package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func NestestDecodeRegisterFlag(field string) byte {
	tokens := strings.Split(field, ":")
	_ = tokens

	result, err := hex.DecodeString(tokens[1])
	if err != nil {
		panic(err)
	}
	if len(result) < 1 {
		panic(fmt.Errorf("error decoding hex to string: %s", field))
	}

	return result[0]
}
