package main

import (
	"fmt"

	"github.com/gorilla/schema"
)

func main() {
	basicModelBinding()
}

var decoder *schema.Decoder

func basicModelBinding() {
	form := map[string][]string{
		"FirstName": []string{"Tony"},
		"LastName":  []string{"Kung"},
	}

	var p Person

	if decoder == nil {
		decoder = schema.NewDecoder()
	}

	decoder.Decode(&p, form)
	fmt.Println(p)
}

type Person struct {
	FirstName string
	LastName  string
}
