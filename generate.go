package randompassword

import (
	"crypto/rand"
	"errors"
	"math/big"
	"slices"
)

var lower_case = []byte{
	'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x',
	'y', 'z',
}

var upper_case = []byte{
	'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z',
}

var number_case = []byte{
	'1', '2', '3', '4', '5', '6', '7', '9', '0',
}

var symbol_case = []byte{
	'~', '`', '!', '@', '#', '$', '%',
	'^', '&', '*', '(', ')', '-', '=',
	'_', '+', '\\', ']', '[', '{', '}',
	'|', ';', '\'', ':', '"', ',', '.',
	'/', '<', '>', '?',
}
var default_options []Options = []Options{LowerCase, Number}

func getOptionsValues(options ...Options) []byte {
	result := make([]byte, 0, 100)
	for _, v := range options {
		switch v {
		case LowerCase:
			result = append(result, lower_case...)
		case UpperCase:
			result = append(result, upper_case...)
		case Number:
			result = append(result, number_case...)
		case Symbol:
			result = append(result, symbol_case...)
		default:
			result = append(result, getOptionsValues(default_options...)...)
		}
	}
	return result
}

func RandomPassword(length uint, options ...Options) (string, error) {
	if len(options) == 0 {
		options = append(options, default_options...)
	}

	if length > 50 {
		return "", errors.New("length to long")
	}

	charset := getOptionsValues(options...)
	length_charset := big.NewInt(int64(len(charset)))

	result := make([]byte, length)

	for i := 0; i < int(length); i++ {
		n, err := rand.Int(rand.Reader, length_charset)
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Uint64()]
	}

	return string(result), nil
}

func RandomPasswordUnique(length uint, options ...Options) (string, error) {
	if len(options) == 0 {
		options = append(options, default_options...)
	}
	if length > 50 {
		return "", errors.New("length to long")
	}
	charset := getOptionsValues(options...)
	length_charset := big.NewInt(int64(len(charset)))
	result := make([]byte, length)

	retry := 0
	success := 0
	for {
		if success >= int(length) || retry == int(length) {
			break
		}
		n, err := rand.Int(rand.Reader, length_charset)
		if err != nil {
			return "", err
		}
		if slices.Contains(result, charset[n.Uint64()]) {
			retry += 1
			continue
		}
		result[success] = charset[n.Uint64()]
		success += 1
		retry -= 1
	}
	return string(result), nil
}
