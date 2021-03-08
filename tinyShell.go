// Implements a tiny shell that we can use to run in our Linux container
// if we don't want to import other Linux binaries. It emulates a few
// basic Linux commands: `cat`, `cd`, `env`, `hostname`, `ls`, and `pwd`.
//
// Compile:
//		$ GOOS=linux go build -o ts tinyShell.go
// Run:
//		$ ./ts
//
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Tiny shell started")
	fmt.Println("Valid commands: cat, cd [dir], env, hostname, ls, pwd, quit")
	pwd, _ := filepath.Abs(".")
	home := pwd
	for true {
		cmd, arg := readCommand("ts: ")
		if cmd == "quit" || cmd == "exit" {
			break
		}

		switch {
		case cmd == "cat":
			cat(pwd, arg)
		case cmd == "env":
			env()
		case cmd == "hostname":
			hostname()
		case cmd == "ls":
			ls(pwd, arg)
		case cmd == "pwd":
			fmt.Printf("%s\n", pwd)
		case cmd == "cd":
			if arg == "" {
				pwd = cd(pwd, home)
			} else {
				pwd = cd(pwd, arg)
			}
		case cmd == "":
			// nothing to do
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
	fmt.Println("Tiny shell ended")
}

func cat(pwd string, filename string) {
	fullname, _ := filepath.Abs(filepath.Join(pwd, filename))
	bytes, err := ioutil.ReadFile(fullname)
	if err == nil {
		fmt.Printf("%s", string(bytes))
	} else {
		fmt.Printf("Error reading %s: %s\n", fullname, err)
	}
}

func cd(pwd string, dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	newPwd, _ := filepath.Abs(filepath.Join(pwd, dir))
	return newPwd
}

func env() {
	for _, env := range os.Environ() {
		fmt.Printf("%s\n", env)
	}
}

func hostname() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Hostname: %s\n", hostname)
	}
}

func ls(pwd string, dir string) {
	var path string
	if filepath.IsAbs(dir) {
		path = dir
	} else {
		path, _ = filepath.Abs(filepath.Join(pwd, dir))
	}

	fmt.Printf("Files in: %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	for _, f := range files {
		fmt.Println("\t" + f.Name())
	}
}

func readCommand(prompt string) (string, string) {
	fmt.Printf("%s", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	tokens := strings.Split(strings.TrimSpace(text), " ")
	cmd := tokens[0]
	arg := ""
	if len(tokens) == 2 {
		arg = tokens[1]
	}
	return cmd, arg
}
