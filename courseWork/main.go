package main

import "log"

var s = "seven"

func main() {

	log.Println(s)
}

func saySomething(s string) (string, string) {
	return s, "world"

}
