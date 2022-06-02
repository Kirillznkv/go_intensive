package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func getArgs() []string {
	var res []string
	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func main() {
	if len(os.Args) > 1 {
		cmdName := os.Args[1]
		var args []string
		args = append(args, os.Args[2:]...)
		args = append(args, getArgs()...)
		cmd := exec.Command(cmdName, args...)
		out, _ := cmd.Output()
		fmt.Printf("%s", string(out))
	}
}
