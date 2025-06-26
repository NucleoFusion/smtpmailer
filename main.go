package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Mail struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func main() {
	godotenv.Load(".env")

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var mail Mail
		if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		Send(mail.Name, mail.Email, mail.Subject, MsgMaker(mail.Name, mail.Email, mail.Content))
	})

	fmt.Println("Listening at " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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

	fmt.Println(msg)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}
