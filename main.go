package main

import (
	"fmt"
)

func main() {
	args := &ClientArgs{
		Host:     "bachelor",
		Port:     "8999",
		Username: "score",
		Password: "***",
		Room:     "bachelor",
		Version:  "1.2.255",
	}

	_, err := NewClient(args)
	if err != nil {
		fmt.Println(err)
	}
	for {
	}
}
