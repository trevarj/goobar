package modules

import (
	"fmt"
	. "goobar/themes"
	"log"
	"os"
	"time"
)

const devpath = "/sys/class/net/enp3s0/operstate"
const icon = "󰈀"

type network struct {
	value string
}

func Network() *network {
	return &network{}
}

func (net *network) Run(chan<- struct{}) {
	f, err := os.Open(devpath)
	if err != nil {
		log.Fatalf("Couldn't open net device file [%s]: %s", devpath, err)
	}
	defer f.Close()

	format := func(s string) string {
		color := Nord.Aurora11
		if s != "up" {
			color = Nord.Aurora14
		}
		return fmt.Sprintf("%%{F%s} %s %%{F-}", color, icon)
	}

	for {
		bytes := make([]byte, 8)
		n, err := f.Read(bytes)
		var val string
		if err != nil {
			val = fmt.Sprintf("net device read error: %s", err)
		} else {
			val = string(bytes[:n])
		}
		net.value = format(val)
		time.Sleep(10 * time.Second)
	}
}

func (n *network) String() string {
	return n.value
}
