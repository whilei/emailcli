package main

import (
	"log"
	"net/smtp"
	"os"
	"strings"

	"io/ioutil"

	"github.com/jordan-wright/email"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	username = kingpin.Flag("username", "Username to authenticate to the SMTP server with").Envar("EMAIL_USERNAME").String()
	password = kingpin.Flag("password", "Password to authenticate to the SMTP server with").Envar("EMAIL_PASSWORD").String()

	//usetls = kingpin.Flag("use-tls", "Use TLS to authenticate").Envar("EMAIL_USETLS").Bool()
	host = kingpin.Flag("host", "Hostname").Envar("EMAIL_HOST").String()
	port = kingpin.Flag("port", "Port number").Envar("EMAIL_PORT").Default("25").Uint16()

	tlsHost = kingpin.Flag("tls-host", "Hostname to use for verifying TLS (default to host if blank)").Envar("EMAIL_TLSHOST").String()

	attachments = kingpin.Flag("attach", "Files to attach to the email.").Envar("EMAIL_ATTACH").ExistingFiles()

	subject = kingpin.Flag("subject", "Subject line of email.").Envar("EMAIL_SUBJECT").String()
	body    = kingpin.Flag("body", "Body of email. Read from stdin if blank.").Envar("EMAIL_BODY").String()
	cc      = kingpin.Flag("cc", "Carbon copy email target. Comma-separated.").Envar("EMAIL_CC").String()

	from = kingpin.Flag("from", "From address for email").Envar("EMAIL_FROM").String()
	to   = kingpin.Arg("to", "Email recipients").Strings()

	timeout  = kingpin.Flag("timeout", "Timeout for mail sending").Envar("EMAIL_TIMEOUT").Duration()
	poolsize = kingpin.Flag("concurrent-sends", "Max concurrent email send jobs").Envar("EMAIL_CONCURRENT_SENDS").Default("1").Int()

	sslInsecure = kingpin.Flag("insecure-skip-verify", "Disable TLS certificate authentication").Envar("EMAIL_INSECURE").Default("false").Bool()
	sslCA       = kingpin.Flag("cacert", "Specify a custom CA certificate to verify against").Envar("EMAIL_CACERT").String()
)

var Version = "0.0.0-dev"

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version(Version)
	kingpin.Parse()

	if *timeout == 0 {
		*timeout = -1
	}

	var bodytxt []byte
	if *body == "" {
		var err error
		bodytxt, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			println(err)
			os.Exit(1)
		}
	} else {
		bodytxt = []byte(*body)
	}

	if *from == "" {
		from = username
	}

	if *username == "" {
		log.Fatal("empty username")
	} else {
		log.Println("username:", *username)
	}

	if *password == "" {
		log.Fatal("empty pass")
	} else {
		// log.Println("pass:", *password)
	}

	var err error
	e := email.NewEmail()

	if *cc != "" {
		e.Cc = strings.Split(*cc, ",")
	}

	e.From = *from
	e.To = *to
	e.Subject = *subject
	e.Text = bodytxt
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", *username, *password, "smtp.gmail.com"))

	if err != nil {
		println("Error sending mail:", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
