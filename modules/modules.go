package modules

type ModResult struct {
	From   string
	Result string
}

type Module interface {
	Run(chan<- struct{})
}
