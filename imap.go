package main

import (
	"bufio"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
	"strings"
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
		/*
			to, err := header.AddressList("To")
			if err != nil {
				log.Fatal(err)
			}
		*/
		subject, err := header.Subject()
		if err != nil {
			log.Fatal(err)
		}
		if from[0].Address == "noreply@dtf.ru" && subject == "Подтверждение регистрации" {
			rawMessage := readMessageParts(messageReader)
			if rawMessage == "" {
				break
			}
			scanner := bufio.NewScanner(strings.NewReader(rawMessage))
			for scanner.Scan() {
				if strings.HasPrefix(scanner.Text(), "[Подтвердить]") {
					link := strings.TrimLeft(strings.TrimRight(scanner.Text(), ")"), "[Подтвердить](")
					fmt.Println(link)

				}
			}
		}
	}
}

func readMessageParts(msgReader *mail.Reader) string {
	for {
		p, err := msgReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case mail.TextHeader:
			b, _ := ioutil.ReadAll(p.Body)
			return string(b)
		case mail.AttachmentHeader:
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}
	}
	return ""
}
