package snyk_iac_capture

type Command struct {
	Org string
}

func (c Command) Run() int {
	return 0
}
