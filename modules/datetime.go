package modules

import "time"

type datetime struct {
	ticker time.Ticker
	value  string
}

var format = "Mon Jan 1 3:04 PM"

func DateTime() *datetime {
	return &datetime{
		value: "",
	}
}

func (dt *datetime) Run(updateChannel chan<- struct{}) {
	var minute int
	for {
		now := time.Now()
		dt.value = now.Format(format)
		if now.Minute() != minute {
			minute = now.Minute()
			updateChannel <- struct{}{}
		}
		time.Sleep(time.Second)
	}

}

func (dt *datetime) String() string {
	return dt.value
}
