package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/gosuri/uilive"
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
	link      string
	ip        string
	status    status
	lastCheck time.Time
}

func main() {
	applications := []*application{
		&application{
			link:   "kromline.ee",
			status: INIT,
		},
		&application{
			link:   "admin.kromline.ee",
			status: INIT,
		},
		&application{
			link:   "alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "jenkins.alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "docker.alekseikromski.com",
			status: INIT,
		},
		&application{
			link:   "blog.alekseikromski.com",
			status: INIT,
		},
	}
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	go func() {
		for {
			for _, app := range applications {
				pinger, err := ping.NewPinger(app.link)
				if err != nil {
					panic(err)
				}

				pinger.Count = 1
				pinger.Timeout = 2 * time.Second
				pinger.Run() // blocks until finished
				stats := pinger.Statistics()

				if stats.PacketLoss > 0.0 {
					app.status = FAILED
				} else {
					app.status = SUCCESS
				}
				app.ip = stats.IPAddr.String()
				app.lastCheck = time.Now()
				time.Sleep(1 * time.Second)
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
		content += fmt.Sprintf("[%s]: %s [%s] (%.fs ago)\n", status, app.link, app.ip, time.Now().Sub(app.lastCheck).Seconds())
	}
	fmt.Fprintf(writer, content)

}
