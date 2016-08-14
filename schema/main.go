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
		"firstName":         []string{"Tony"},
		"lastName":          []string{"Kung"},
		"addresses.0.city":  []string{"Hong Kong"},
		"addresses.0.state": []string{"N/A"},
		"addresses.1.city":  []string{"Hong Kong 2"},
		"addresses.1.state": []string{"N/A 2"},
	}

	var p Person

	if decoder == nil {
		decoder = schema.NewDecoder()
	}

	decoder.Decode(&p, form)
	fmt.Println(p)
}

type Person struct {
	FirstName string    `schema:"firstName"`
	LastName  string    `schema:"-"` //skip this field
	Addresses []Address `schema:"addresses"`
}

type Address struct {
	City  string `schema:"city"`
	State string `schema:"state"`
}
