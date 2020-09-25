package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
        "text/template"
	"time"
)

type Configuration struct {
	Sourcefile                  string
	Logfile                     string
	Template                    string
	Sleep_time_in_hours         int
	Debug                       int
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

func main(){
	fmt.Println(techniques[rand.Intn(num_techniques)])

		// read template
	tmpl, err := template.ParseFiles(configuration.Template)
        check(err)
	
	   data := TemplateData{
            Technique: techniques[rand.Intn(num_techniques)],         
	   }
	
        tmpl.Execute(os.Stdout, data)
	logger.Print("Ended...")
}
