package main

import (
	"flag"
	"fmt"

	com "./compare_db"
)

func main() {
	oldName := flag.String("old", "", "old file name")
	newName := flag.String("new", "", "new file name")
	flag.Parse()
	if *oldName == "" || *newName == "" {
		fmt.Println("ERROR: file name")
	} else {
		com.Compare(*oldName, *newName)
	}
}
