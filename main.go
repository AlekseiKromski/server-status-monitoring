package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"net/http"
	"time"
)

type status string

const (
	FAILED  status = "FAILED"
	SUCCESS status = "SUCCESS"
	INIT    status = "INIT"
)

var (
	info    = color.New(color.FgWhite, color.BgCyan)
	success = color.New(color.FgWhite, color.BgGreen)
	failed  = color.New(color.FgWhite, color.BgRed)
)

type application struct {
	link       string
	status     status
	statusCode int
	lastCheck  time.Time
}

func main() {
	applications := []*application{
		&application{
			link:   "https://kromline.ee",
			status: INIT,
		},
		&application{
			link:   "https://admin.kromline.ee",
			status: INIT,
		},
		&application{
			link:   "https://alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "https://jenkins.alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "https://docker.alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "https://blog.alekseikromski.com",
			status: INIT,
		},
	}
	successStatusCodes := []int{200, 400, 500, 403}

	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	go func() {
		for {
			for _, app := range applications {
				req, err := http.NewRequest(http.MethodGet, app.link, nil)
				if err != nil {
					app.status = FAILED
					continue
				}
				res, err := http.DefaultClient.Do(req)
				app.statusCode = res.StatusCode
				app.lastCheck = time.Now()

				if contains(res.StatusCode, successStatusCodes) && err == nil {
					app.status = SUCCESS
				} else {

					app.status = FAILED
					continue
				}

			}
			time.Sleep(5 * time.Second)
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

func render(applications []*application, writer *uilive.Writer) {
	content := ""
	for _, app := range applications {
		status := ""
		switch app.status {
		case INIT:
			status = info.Sprintf("%s", app.status)
		case SUCCESS:
			status = success.Sprintf("%s", app.status)
		case FAILED:
			status = failed.Sprintf("%s", app.status)
		}
		content += fmt.Sprintf("[%s][%d]: %s (%.fs ago)\n", status, app.statusCode, app.link, time.Now().Sub(app.lastCheck).Seconds())
	}
	fmt.Fprintf(writer, content)

}
