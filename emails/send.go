package main

import (
    "log"
    "net/smtp"
	"time"
)

func main() {
    send("hello bro, This just test.")
}

func send(body string) {
    from := "znaumets@gmail.com"
    pass := "dydechqszncoqhck"
    to := "kchkskpr@gmail.com"

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: Hello there 2\n\n" +
        body
	
	start_time := time.Now()
	log.Println("Start sending email at: " + start_time.String())

    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
        from, 
		[]string{to}, 
		[]byte(msg),
	)

	end_time := time.Now()
	log.Println("End sending email at: " + end_time.String())

    if err != nil {
        log.Printf("smtp error: %s", err)
        return
    }
    log.Println("Successfully sended to " + to)
}
