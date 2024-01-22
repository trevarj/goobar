package modules

import (
	"fmt"
	. "goobar/themes"
	"log"
	"net"
	"strings"
	s "strings"
)

const (
	workspace         = ""
	workspace_focused = ""
)

type bspwm struct {
	name  string
	state string
}

func Bspwm() *bspwm {
	return &bspwm{name: "bspwm"}
}

func formatted(report string) string {
	var formatted string = ""
	splits := s.Split(report, ":")
	log.Println(splits)
	for _, s := range splits[1:] {
		var name, icon, fg, bg string
		// REPORT FORMAT from man bspc
		switch s[0] {
		case 'O': // occupied focused
			fg = Nord.Aurora13
			bg = Nord.PolarNight3
			icon = workspace_focused
		case 'o': // occupied unfocused
			fg = Nord.Aurora13
			bg = Nord.PolarNight2
			icon = workspace
		case 'F': // free focused
			fg = Nord.Aurora13
			bg = Nord.PolarNight3
			icon = workspace_focused
		case 'f': // free unfocused
			fg = Nord.Aurora13
			bg = Nord.PolarNight2
			icon = workspace
		case 'U': // urgent focused
			fg = Nord.Aurora11
			bg = Nord.PolarNight3
			icon = workspace_focused
		case 'u': // urgent unfocused
			fg = Nord.Aurora11
			bg = Nord.PolarNight2
			icon = workspace
		}
		if icon != "" {
			name = s[1:]
			formatted += fmt.Sprintf("%%{F%s}%%{B%s}%%{A:bspc desktop -f %s:}  %s  %%{A}%%{B-}%%{F-}%%{-u}", fg, bg, name, icon)
		}
	}
	return formatted
}

func (m *bspwm) Run(updateChannel chan<- struct{}) {
	// TODO: user bspc --print-socket-path on master version
	sock, err := net.Dial("unix", "/tmp/bspwm_0_0-socket")
	defer sock.Close()
	if err != nil {
		log.Fatal("Couldn't connect to bspwm socket")
		return
	}

	_, err = sock.Write([]byte("subscribe\x00"))
	if err != nil {
		log.Fatal("Couldn't subscribe to bspwm report")
	}

	reportChannel := make(chan string)
	for {
		go func() {
			bytes := make([]byte, 128)
			n, err := sock.Read(bytes)
			if err != nil {
				log.Fatal("Couldn't read from bspwm socket")
			}
			// sometimes it spits out multiple lines
			report := strings.Split(string(bytes[:n]), "\n")[0]
			reportChannel <- formatted(report)
		}()
		m.state = <-reportChannel
		updateChannel <- struct{}{}
	}
}

func (m *bspwm) String() string {
	return m.state
}
