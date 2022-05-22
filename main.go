package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PersoSirEduard/kode/kode"
)

const (
	VERSION = "alpha 0.5.0"
)

func main() {

	path := flag.String("run", "main.kd", "Path to the Kode file.")
	showVersion := flag.Bool("version", false, "Show the current version of Kode.")
	StdIn := flag.Bool("runStdIn", false, "Read from stdin.")
	flag.Parse()

	if *StdIn {
		in := bufio.NewScanner(os.Stdin)

		code := ""

		for in.Scan() {

			bytes := in.Bytes()
			txt := string(bytes)

			if strings.ReplaceAll(txt, " ", "") == "exit" {
				break
			}

			if txt != "" {
				code += txt + "\n"
			}

		}

		err := kode.Run(code)

		if err != nil {
			fmt.Println(err.Error())
		}

		os.Exit(1)
	}

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
