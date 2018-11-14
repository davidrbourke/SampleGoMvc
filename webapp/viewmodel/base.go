package viewmodel

// Base is the viewmodel base struct
type Base struct {
	Title string
}

// NewBase returns constructed Base
func NewBase() Base {
	return Base{
		Title: "Lemonade",
	}
}
