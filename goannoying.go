// Wrote on the fly reading
// http://nathanleclaire.com/blog/2013/12/17/sending-email-from-gmail-using-golang/
// https://code.google.com/p/go-wiki/wiki/SendingMail
// to show how to write annoying mail to your supervisor
// to have an answer or to have a paper read!:P
package main

import (
	"bytes"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"text/template"
	"time"

	"code.google.com/p/gopass"
)

func main() {
	var (
		body         bytes.Buffer
		err          error
		mailTemplate *template.Template
		pw           string
		randomSeed   int64 = 42
	)

	if pw, err = gopass.GetPass("Please insert yuor password: "); err != nil {
		log.Fatal("Error retrieving password; ", err)
	}

	log.Println("Creating user data and generating mail text")
	userData := &UserData{
		From:        "",
		Password:    pw,
		EmailServer: "smtp.gmail.com",
		Port:        587,
		To:          []string{"", ""},
		Subject:     "Come scrivere una mail automatica in go",
		Body: `Esempio di mail automatica mandata a tempi variabili 
per stressare il supervisor e ottenere una risposta!:P
		
		by brunetto
		`,
		MaxWaitingTime: 10,
		NMails:         5,
	}

	log.Println("Creating authorization")
	auth := smtp.PlainAuth("",
		userData.From,
		userData.Password,
		userData.EmailServer,
	)

	log.Println("Filling mail template")
	mailTemplate = template.New("emailTemplate")
	if mailTemplate, err = mailTemplate.Parse(EmailTemplate); err != nil {
		log.Fatal("error trying to parse mail template")
	}

	err = mailTemplate.Execute(&body, userData)
	if err != nil {
		log.Fatal("error trying to execute mail template")
	}

	log.Println("Init random seed")
	rand.Seed(randomSeed)

	log.Println("Start loop")
	for idx := 0; idx < userData.NMails; idx++ {
		log.Println("Sending mail #", idx+1)
		err = smtp.SendMail(userData.EmailServer+":"+strconv.Itoa(userData.Port), // in our case, "smtp.google.com:587"
			auth,
			userData.From,
			userData.To,
			body.Bytes())
		if err != nil {
			log.Fatal("ERROR: attempting to send a mail ", err)
		} else {
			log.Println("Mail sent succesfully!:)")
		}
		waitingTime := time.Duration(rand.Intn(userData.MaxWaitingTime)) * time.Minute
		log.Println("Waiting ", waitingTime)
		time.Sleep(waitingTime)
	}

	log.Println("Done")
}

type UserData struct {
	From           string
	Password       string
	EmailServer    string
	Port           int
	To             []string
	Subject        string
	Body           string
	MaxWaitingTime int
	NMails         int
}

const EmailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`
