package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/PersoSirEduard/kode/kode"
)

const (
	VERSION = "alpha 0.5.0"
)

func main() {

	path := flag.String("run", "main.kd", "Path to the Kode file.")
	showVersion := flag.Bool("version", false, "Show the current version of Kode.")
	flag.Parse()

	if *showVersion {
		fmt.Println("Kode version: " + VERSION)
		fmt.Println("Created by Eduard A. and Frédéric J.")
		fmt.Println("Make sure to check out the GitHub repository: https://github.com/PersoSirEduard/Kode")
		return
	}

	code, err := ioutil.ReadFile(*path)

	if err != nil {
		println("Error: Could not find and read the file \"" + *path + "\".")
	}

	err = kode.Run(string(code))

	if err != nil {
		println(err.Error())
	}

}
