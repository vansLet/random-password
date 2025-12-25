package randompassword

type Options int

const (
	LowerCase Options = iota
	UpperCase
	Number
	Symbol
)

var nameOptions []string = []string{
	"LowerCase",
	"UpperCase",
	"Number",
	"Symbol",
}

func (option Options) String() string {
	if option >= LowerCase && option <= Symbol {
		return nameOptions[option]
	}
	return "None"
}
