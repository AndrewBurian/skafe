package main

import (
	"fmt"
	audit "github.com/mozilla/libaudit-go"
)

func main() {
	fmt.Println("Hello SKAFE!")
	s, err := audit.NewNetlinkConnection()

	if err != nil {
		fmt.Println(err)

	}

	defer s.Close()

}
