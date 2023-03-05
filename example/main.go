package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/taniko/gothic/event"
)

func main() {
	f, err := os.Open("model/main.smithy")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	events, err := event.Parse(bufio.NewScanner(f))
	if err != nil {
		panic(err)
	}
	for events.Next() {
		e := events.Event()
		fmt.Println(e.Type(), e.Value())
	}
}
