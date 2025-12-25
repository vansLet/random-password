package randompassword

import (
	"crypto/rand"
	"errors"
	"math/big"
	"slices"
)

//# base-value

// base value lower case
var lower_case = []byte{
	'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x',
	'y', 'z',
}

// base value upper case
var upper_case = []byte{
	'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z',
}

// base value number
var number_case = []byte{
	'1', '2', '3', '4', '5', '6', '7', '9', '0',
}

// base value symbol case
var symbol_case = []byte{
	'~', '`', '!', '@', '#', '$', '%',
	'^', '&', '*', '(', ')', '-', '=',
	'_', '+', '\\', ']', '[', '{', '}',
	'|', ';', '\'', ':', '"', ',', '.',
	'/', '<', '>', '?',
}

//#

// default value options
var default_options []Options = []Options{LowerCase, Number}

// mendapatkan slice byte yang berisi base-value yang sudah digabungkan.
//
// penggabungan base-value ditentukan oleh parameter params options,
// jika paramater options diisi bukan dari type Options akan menyebabkan panic
//
// pastikan params harus unique, karena func ini tidak melakukan pengecekan.
func getOptionsValues(options ...Options) []byte {
	result := make([]byte, 0, 100)
	for _, v := range options {
		switch v { // <- melakukan pengecekan Options, yang menentukan base-value ditambkan ke var result
		case LowerCase:
			result = append(result, lower_case...)
		case UpperCase:
			result = append(result, upper_case...)
		case Number:
			result = append(result, number_case...)
		case Symbol:
			result = append(result, symbol_case...)
		default:
			panic("required type Options") // <- jika bukan dari type Options langsung menyebabkan panic
		}
	}
	return result
}

// mengembalikan slices Options yang valuenya unique.
//
// jika paramater options diisi bukan dari type Options akan menyebabkan panic
func sortingUniqueOptions(options []Options) []Options {
	if len(options) == 1 { // <- langsung di return
		return options
	}
	result := make([]Options, 0, 4)
	for _, o := range options {
		switch o {
		case LowerCase, Number, Symbol, UpperCase:
			{
				// jika o ada di result di skip, jika tidak ditambahkan ke result
				if slices.Contains(result, o) {
					continue
				}
				result = append(result, o)
			}
		default:
			panic("required type Options") // <- jika bukan dari type Options langsung menyebabkan panic
		}
	}
	return result
}

// mengembalikan sebuah string random-password.
//
// parameter length: set panjang random-password, jika lebih dari 50 mengembalikan error.
//
// parameter params options: set Option yang akan di generated, jika kosong secara default{ LowerCase, Number }.
//
// menggunakan crypto/rand yang dimana secara bawa'an sudah thread safe
func RandomPassword(length uint, options ...Options) (string, error) {
	if len(options) == 0 { // <- pengecekan jika options kosong akan di isi default_options,jika tidak akan melakukan sortingUnique
		options = default_options
	} else {
		options = sortingUniqueOptions(options)
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

// mengembalikan sebuah string random-password yang bersifat unique.
//
// parameter length: set panjang random-password, jika lebih dari 50 mengembalikan error.
//
// parameter params options: set Option yang akan di generated, jika kosong secara default{ LowerCase, Number }.
//
// terkadang panjang dari result tidak sama dengan parameter length,
// menggunakan crypto/rand yang dimana secara bawa'an sudah thread safe
func RandomPasswordUnique(length uint, options ...Options) (string, error) {
	if len(options) == 0 { // <- pengecekan jika options kosong akan di isi default_options, jika tidak akan melakukan sortingUnique
		options = default_options
	} else {
		options = sortingUniqueOptions(options)
	}
	if length > 50 {
		return "", errors.New("length to long")
	}

	charset := getOptionsValues(options...)
	length_charset := big.NewInt(int64(len(charset)))
	result := make([]byte, length)

	retry := 0   //jumlah data yang gagal di tambahkan
	success := 0 // jumlah data yang berhasil di tambahkan
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
