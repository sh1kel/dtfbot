package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"log"
	"sync"
)

func listImapMessages(config *Configuration, messages chan *imap.Message, imapClient *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	defer imapClient.Logout()
	err := imapClient.Login(config.MailServer.User, config.MailServer.Password)
	if err != nil {
		log.Fatal(err)
	}

	mailBox, err := imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mailBox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}

	fmt.Printf("Number of messages: %d\n", mailBox.Messages)
	from := uint32(1)
	to := mailBox.Messages
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	err = imapClient.Fetch(seqSet, items, messages)
	if err != nil {
		log.Fatal(err)
	}

}

func parseMessage(messages chan *imap.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	section := &imap.BodySectionName{}
	for msg := range messages {
		if msg == nil {
			log.Fatal("Can't read message from server")
		}

		msgBody := msg.GetBody(section)
		if msgBody == nil {
			log.Fatal("Can't read message body from server")
		}

		messageReader, err := mail.CreateReader(msgBody)
		if err != nil {
			log.Fatal(err)
		}
		header := messageReader.Header
		from, err := header.AddressList("From")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("From: %s\n", from)
		subject, err := header.Subject()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Subject: %s\n", subject)
	}
}
