package main

import (
	"io/ioutil"

	"github.com/PersoSirEduard/kode/kode"
)

func main() {

	code, err := ioutil.ReadFile("program.k")

	if err != nil {
		panic(err)
	}

	err = kode.Run(string(code))

	if err != nil {
		panic(err)
	}

}
