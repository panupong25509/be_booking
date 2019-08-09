package mailer

import (
	"log"

	"github.com/panupong25509/be_booking_sign/config"

	"gopkg.in/gomail.v2"
)

func HTML(status, booking string, comment string) string {
	return `
<html>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>  
<body style="background-color: #ECF3F6">
    <div style="padding: 50px 100px;">
      <div style="background-color: white;padding: 30px 50px;">
        <div style="border-bottom: 2px solid #ECF3F6">
          <h1>Your book ` + status + `</h1>
        </div>

        <p>You can check status of ` + booking + ` on dashboard or click to button</p>
				<p>` + comment + `</p>
				<a href="http://localhost:3001/">
          <button
            style="background-color: #f47836;color: white;border: none;padding: 8px 10px;border-radius: 8px;cursor: pointer;"
          >
            Click here
          </button>
        </a>
      </div>
    </div>
  </body>
</html>
`
}

func SendEmail(Subject string, email string, status string, comment string) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.GetConfig().MAIL_EMAIL)
	m.SetHeader("To", email)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", Subject)
	html := HTML(status, "booking", comment)
	m.SetBody("text/html", html)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GetConfig().MAIL_EMAIL, config.GetConfig().MAIL_PASSWORD)
	log.Print(config.GetConfig().MAIL_EMAIL)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return
}
