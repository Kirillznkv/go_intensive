package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func wc_w(filename string) {
	var res uint = 0
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res += uint(len(strings.Split(scanner.Text(), " ")))
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d\t%s\n", res, filename)
	i++
}

func wc_l(filename string) {
	var res uint = 0
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res++
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d\t%s\n", res, filename)
	i++
}

func wc_m(filename string) {
	var res uint = 0
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res += uint(len(scanner.Text()))
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d\t%s\n", res, filename)
	i++
}

var i int = 0

func main() {
	w := flag.Bool("w", true, "counting words")
	l := flag.Bool("l", false, "counting lines")
	m := flag.Bool("m", false, "counting characters")
	flag.Parse()

	files := flag.Args()
	for _, f := range files {
		if *l {
			go wc_l(f)
		} else if *m {
			go wc_m(f)
		} else if *w {
			go wc_w(f)
		}
	}
	for i != len(files) {
	}
}
