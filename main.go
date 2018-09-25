package main

import (
	"fmt"
	"github.com/emersion/go-imap"
)

func main() {
	var config Configuration
	mailMessagesChan := make(chan *imap.Message, 100)
	config = loadConfig()
	fmt.Printf("%s\n", config)
	go readImapMessages(&config, mailMessagesChan)

	for msg := range mailMessagesChan {
		fmt.Printf("# %s\n", msg.Envelope.Subject)
	}
}
