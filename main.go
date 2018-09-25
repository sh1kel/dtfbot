package main

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"sync"
)

func main() {
	var config Configuration
	var imapClient *client.Client
	var wg sync.WaitGroup

	mailMessagesChan := make(chan *imap.Message, 100)
	config = loadConfig()

	imapClient, err := client.DialTLS(config.MailServer.Host+":"+config.MailServer.ImapPort, nil)
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(2)
	go listImapMessages(&config, mailMessagesChan, imapClient, &wg)
	go parseMessage(mailMessagesChan, &wg)
	wg.Wait()
}
