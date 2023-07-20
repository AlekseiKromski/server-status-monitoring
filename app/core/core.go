package core

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	gomail "gopkg.in/mail.v2"
	"net/http"
	"strings"
	"time"
)

type status string

const (
	FAILED  status = "FAILED"
	SUCCESS status = "SUCCESS"
	INIT    status = "INIT"
	RETRY   status = "RETRY"
)

var (
	info    = color.New(color.FgWhite, color.BgCyan)
	success = color.New(color.FgWhite, color.BgGreen)
	failed  = color.New(color.FgWhite, color.BgRed)
)

type Application struct {
	Link       string    `json:"link"`
	Status     status    `json:"status"`
	StatusCode int       `json:"statusCode"`
	Sent       bool      `json:"sent"`
	LastCheck  time.Time `json:"lastCheck"`
}

func Start(update chan []*Application, config *Config) {
	applications := []*Application{
		&Application{
			Link:   "https://kromline.alekseikromski.com",
			Status: INIT,
		},
		&Application{
			Link:   "https://kromline-admin.alekseikromski.com",
			Status: INIT,
		},
		&Application{
			Link:   "https://stopper.vaheta.me/",
			Status: INIT,
		},
		&Application{
			Link:   "https://vaheta.me",
			Status: INIT,
		},
		&Application{
			Link:   "https://kromline.ee",
			Status: INIT,
		},
		&Application{
			Link:   "https://admin.kromline.ee",
			Status: INIT,
		},
		&Application{
			Link:   "https://alekseikromski.com",
			Status: INIT,
		},
		&Application{
			Link:   "https://jenkins.alekseikromski.com",
			Status: INIT,
		},
		&Application{
			Link:   "https://docker.alekseikromski.com",
			Status: INIT,
		},
		&Application{
			Link:   "https://blog.alekseikromski.com",
			Status: INIT,
		},
	}
	successStatusCodes := []int{200, 400, 500, 403}
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	//Retry policy
	rp := NewRetryPolicy(config.Count, config.Timeout*time.Second)

	go func() {
		for {
			for _, app := range applications {
				var resp *http.Response
				for rp.Retry() {
					app.Status = RETRY
					res, err := client.Get(app.Link)
					if err != nil {
						rp.Timeout()
						continue
					}

					resp = res
					break
				}

				rp.Clean()

				if resp == nil {
					app.Status = FAILED
					app.StatusCode = 0
					app.LastCheck = time.Now()
					sendEmail(app, config)
					app.Sent = true
					continue
				}

				app.StatusCode = resp.StatusCode
				app.LastCheck = time.Now()

				if contains(resp.StatusCode, successStatusCodes) {
					app.Sent = false
					app.Status = SUCCESS
				} else {
					app.Status = FAILED
					sendEmail(app, config)
					app.Sent = true
					continue
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()

	go func() {
		//update websocket
		for {
			update <- applications
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for {
		render(applications, writer)
		time.Sleep(500 * time.Millisecond)
	}

	writer.Stop() // flush and stop rendering
}

func contains(code int, list []int) bool {
	for _, item := range list {
		if code == item {
			return true
		}
	}
	return false
}

func render(applications []*Application, writer *uilive.Writer) {
	content := ""
	for _, app := range applications {
		status := ""
		switch app.Status {
		case INIT:
			status = info.Sprintf("%s", app.Status)
		case SUCCESS:
			status = success.Sprintf("%s", app.Status)
		case FAILED:
			status = failed.Sprintf("%s", app.Status)
		case RETRY:
			status = info.Sprintf("%s", app.Status)
		}
		content += fmt.Sprintf("[%s][%d]: %s (%.fs ago)\n", status, app.StatusCode, app.Link, time.Now().Sub(app.LastCheck).Seconds())
	}
	fmt.Fprintf(writer, content)

}

func sendEmail(app *Application, config *Config) {
	if app.Sent {
		return
	}

	app.Sent = true
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", config.EmailUsername)

	// Set E-Mail receivers
	m.SetHeader("To", config.EmailTo)

	// Set E-Mail subject
	m.SetHeader("Subject", "Failed")

	// Set E-Mail body. You can set plain text or html with text/html
	status := ""
	switch app.Status {
	case INIT:
		status = fmt.Sprintf("%s", app.Status)
	case SUCCESS:
		status = fmt.Sprintf("%s", app.Status)
	case FAILED:
		status = fmt.Sprintf("%s", app.Status)
	}
	_, address, _ := strings.Cut(app.Link, "https://")
	m.SetBody("text/plain", fmt.Sprintf("[%s][%d]: %s (%s)\n", status, app.StatusCode, address, app.LastCheck))

	// Settings for SMTP server
	d := gomail.NewDialer(config.EmailHost, config.EmailPort, config.EmailUsername, config.EmailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	d.DialAndSend(m)
}
