package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"
)

type Configuration struct {
	Sourcefile          string
	Logfile             string
	Template            string
	Sleep_time_in_hours int
	Mailhost            string
	Mailport            int
	Mailuser            string
	Mailpwd             string
	Email               string
	Debug               int
}

type TemplateData struct {
	Technique string
}

var logfile *os.File
var err error
var logger *log.Logger
var techniques []string
var num_techniques int
var configuration Configuration

// init opens a log file, reads the techniques file
// and creates a new random seed
func init() {

	// config
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	check(err)

	// logging
	logfile, err = os.OpenFile(configuration.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check(err)
	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Print("Started...")

	// init random generator
	rand.Seed(time.Now().UnixNano())

	// read source file
	content, err := ioutil.ReadFile(configuration.Sourcefile)
	check(err)
	techniques = strings.Split(string(content), "\n")

	num_techniques = len(techniques) - 1
	logger.Printf("Found %d elements in %s\n", num_techniques, configuration.Sourcefile)

}

// check panics if an error is detected
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	technique := techniques[rand.Intn(num_techniques)]

	// read template
	tmpl, err := template.ParseFiles(configuration.Template)
	check(err)

	data := TemplateData{
		Technique: technique,
	}

	//sending email
	var body bytes.Buffer
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n"
	subject := fmt.Sprintf("Subject: Aikido-Erinnerung: %s\n", technique)
	to := fmt.Sprintf("To: %s\n", configuration.Email)
	from := fmt.Sprintf("From: %s\n", configuration.Email)

	body.Write([]byte(mime + from + to + subject + "\n\n"))
	posteoAuth := smtp.PlainAuth("", configuration.Mailuser, configuration.Mailpwd, configuration.Mailhost)

	tmpl.Execute(&body, data)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", configuration.Mailhost, configuration.Mailport), posteoAuth, configuration.Mailuser, []string{configuration.Email}, body.Bytes())

	if err != nil {
		logger.Printf("Error SendMail: %s", err)
	} else {
		logger.Printf("Send email to '%s' with technique '%s'", configuration.Email, technique)
	}
	logger.Print("Ended...")
}
