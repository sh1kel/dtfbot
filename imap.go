package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

func readImapMessages(config *Configuration, messages chan *imap.Message) {
	imapClient, err := client.DialTLS(config.MailServer.Host+":"+config.MailServer.ImapPort, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer imapClient.Logout()
	err = imapClient.Login(config.MailServer.User, config.MailServer.Password)
	if err != nil {
		log.Fatal(err)
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- imapClient.List("", "*", mailboxes)
	}()

	mailBox, err := imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Flags for INBOX:", mailBox.Flags)

	from := uint32(1)
	to := mailBox.Messages
	fmt.Printf("Number of messages: %d\n", to)
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	done = make(chan error, 1)
	go func() {
		done <- imapClient.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}
