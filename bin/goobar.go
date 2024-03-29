package main

import (
	"fmt"
	. "goobar/modules"
	. "goobar/themes"
	"log"
	"os/exec"
	"strings"
)

var lemonArgs = []string{
	"-n", "bspwm_panel",
	"-g", "x38",
	"-a", "10",
	"-u", "0",
	"-f", "Terminus:style=Bold:size=12",
	"-o", "-10",
	"-f", "JetBrainsMono NFM:size=18",
	"-o", "0",
	"-f", "cryptofont:style=Regular:size=10",
	"-o", "5",
	"-F", Nord.Aurora13,
	"-B", Nord.PolarNight1,
	"-U", Nord.Aurora12,
}

func main() {
	// Create the modules
	modules := map[string]Module{
		"bspwm":    Bspwm(),
		"system":   System(),
		"datetime": DateTime(),
		"network":  Network(),
	}

	// Create a channel to signal an update
	updateChannel := make(chan struct{})

	// Start modules
	for _, module := range modules {
		// create new instance for go routine
		module := module
		go module.Run(updateChannel)
	}

	lemonbarCmd := exec.Command("lemonbar", lemonArgs...)
	lemonStdin, err := lemonbarCmd.StdinPipe()
	if err != nil {
		log.Fatalf("Couldn't get stdin pipe to lemonbar: %s", err)
	}

	lemonStdout, err := lemonbarCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Couldn't get stdout pipe from lemonbar: %s", err)
	}

	// Start a separate goroutine to handle the latest results
	go func() {
		defer lemonStdin.Close()
		for {
			_ = <-updateChannel
			// var output string
			log.Println("Updated Results:")
			for name, mod := range modules {
				log.Printf("%s: %s\n", name, mod)
			}
			bar := fmt.Sprintf("%%{l} %%{T1}%s%%{T-} %s %%{c}%s %%{r}%%{T2}%s%%{T-}",
				modules["bspwm"],
				modules["system"],
				modules["datetime"],
				modules["network"],
			)
			if _, err = lemonStdin.Write([]byte(bar)); err != nil {
				log.Println("Couldn't write to lemonbar stdin: ", err)
			}
			log.Println("------------------------------")
		}
	}()

	// Handle lemonbar's stdout which will be commands from clicks actions
	//
	// TODO:
	// write sophisticated write handler that can determine which module it came
	// from if i need that for modifying the state of a module somehow
	go func() {
		defer lemonStdout.Close()
		for {
			bytes := make([]byte, 128)
			n, err := lemonStdout.Read(bytes)
			if err != nil {
				log.Fatalf("Couldn't read from lemonbar stdout: %s", err)
			}
			cmdString := string(bytes[:n])
			splits := strings.Split(strings.TrimSpace(cmdString), " ")
			cmd := exec.Command(splits[0], splits[1:]...)
			go func() {
				if _, err := cmd.Output(); err != nil {
					log.Printf("Sub command, %s, failed: %s\n", splits, err)
				}
			}()
			log.Println("click action: " + string(bytes))
		}
	}()

	lemonbarCmd.Start()
	// TODO: lock the bar above window manager
	lemonbarCmd.Wait()
}
