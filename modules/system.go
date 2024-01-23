package modules

import (
	"fmt"
	"log"
	"time"

	"goobar/themes"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	cpu_logo  = ""
	cpu_color = themes.Nord.Frost9
	mem_logo  = "󰍛"
	mem_color = themes.Nord.Frost9
)

type system struct {
	value string
}

func System() *system {
	return &system{}
}

func (s *system) Run(updateChannel chan<- struct{}) {
	format := func(cpu int, mem int) string {
		cpu_color := cpu_color
		mem_color := mem_color
		if cpu == -1 {
			cpu_color = themes.Nord.Aurora11
		}
		if cpu == -1 {
			cpu_color = themes.Nord.Aurora11
		}
		return fmt.Sprintf("%%{F%s}%s %d%% %%{F-} %%{F%s}%s %d%% %%{F-}",
			cpu_color, cpu_logo, cpu,
			mem_color, mem_logo, mem)
	}

	cpu.Percent(time.Millisecond, false)
	for {
		c, err := cpu.Percent(0, false)
		cpu_perc := int(c[0])
		if err != nil {
			log.Fatalf("Couldn't get cpu percentage: %s", err)
			cpu_perc = -1
		}
		m, err := mem.VirtualMemory()
		mem_perc := int(m.UsedPercent)
		if err != nil {
			log.Fatalf("Couldn't get cpu percentage: %s", err)
			mem_perc = -1
		}
		s.value = format(cpu_perc, mem_perc)

		updateChannel <- struct{}{}
		time.Sleep(2 * time.Second)
	}
}

func (s *system) String() string {
	return s.value
}
