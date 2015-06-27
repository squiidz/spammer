package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"

	"github.com/go-zoo/bone"
	"github.com/go-zoo/claw"
	mw "github.com/go-zoo/claw/middleware"
)

type User struct {
	ID   string
	Name string
	Mail string
	Send int
}

type Mail struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	MailServer string `json:"mailServer"`
	Port       string `json:"port"`
}

var AUTH smtp.Auth
var CONFIG Config

func init() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &CONFIG)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	AUTH = smtp.PlainAuth("", CONFIG.Username, CONFIG.Password, CONFIG.MailServer)
}

func main() {
	muxx := bone.New()
	clw := claw.New(mw.Logger)

	// API ROUTE
	muxx.Get("/api/:user", clw.Use(func(rw http.ResponseWriter, req *http.Request) {

	}))

	muxx.Post("/api/:user", clw.Use(func(rw http.ResponseWriter, req *http.Request) {
		m := &Mail{}
		err := json.NewDecoder(req.Body).Decode(m)
		if err != nil {
			json.NewEncoder(rw).Encode(err)
			return
		}
		fmt.Printf("Sending mail to %s ...\n", m.To)
		go m.Send()
	}))

	http.ListenAndServe(":8000", muxx)
}

func (m *Mail) Send() error {
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", CONFIG.MailServer, CONFIG.Port),
		AUTH,
		CONFIG.Username,
		[]string{m.To},
		[]byte(m.Body),
	)
	if err != nil {
		return err
	}
	return nil
}
