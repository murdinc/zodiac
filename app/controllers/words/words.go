package main

import (
	"fmt"
	"io/ioutil"
)

func main() {

	fmt.Print("test\n\n")
	files, err := ioutil.ReadDir("./seed")
	if err != nil {
		fmt.Print("Error reading directory")
	} else {
		fmt.Printf("Files %V", files)
	}

}
