package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
	}

	http.HandleFunc("POST /send", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		body := r.PostForm

		name, email, content, err := DecodeBody(&body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		go Send("smthing@gmail.com", "Lapis Nucleo", "TEST", MsgMaker(name, email, content))
	})

	fmt.Println("Listening at " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func DecodeBody(body *url.Values) (string, string, string, error) {
	var (
		name    string
		email   string
		content string
	)

	for key, val := range *body {
		switch key {
		case "email":
			email = val[0]
		case "name":
			name = val[0]
		case "content":
			content = val[0]
		}
	}

	if name == "" || email == "" || content == "" {
		return name, email, content, errors.New("invalid params")
	}

	return name, email, content, nil
}

func MsgMaker(name string, email string, content string) string {
	return "Name: " + name + "\n" + "Email: " + email + "\n" + content
}

func Send(email string, name string, subject string, body string) {
	pass := os.Getenv("GM_PASS")
	to := []string{os.Getenv("GM_1"), os.Getenv("GM_2"), os.Getenv("GM_3")}
	from := os.Getenv("GM_EMAIL")

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	msg := "From: " + email + "\n" + "To: " + to[0] + "\n" + "Subject:" + subject + "\n\n" + body

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}