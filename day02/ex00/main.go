package main

import (
	"flag"
	"log"

	find "./myFind"
)

func main() {
	var fl find.Fl
	fl.F = flag.Bool("f", false, "flag file")
	fl.D = flag.Bool("d", false, "flag directori")
	fl.Sl = flag.Bool("sl", false, "flag symlink")
	fl.Ext = flag.String("ext", "", "flag extension")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Error: file name")
	}
	if *fl.Ext != "" && (*fl.F == false || *fl.D == true || *fl.Sl == true) {
		log.Fatal("Error: flag")
	}
	if *fl.Sl == false && *fl.F == false && *fl.D == false && *fl.Ext == "" {
		*fl.Sl, *fl.F, *fl.D = true, true, true
	}

	dirName := flag.Arg(0)
	find.Find(dirName, &fl)
}
