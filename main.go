package main

import "fmt"

func main() {
	var config Configuration
	config = loadConfig()
	fmt.Printf("%s\n", config)

}
