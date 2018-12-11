package gobacktest

type Cmd struct {
	Event
	cmd string
}

type Result struct {
	Event
	data interface{}
}

func (c *Cmd) SetCmd(cmd string) {
	c.cmd = cmd
}

func (c *Cmd) Cmd() string {
	return c.cmd
}


func (r *Result) SetData(data interface{}) {
	r.data = data
}

func (r *Result) Data() interface{} {
	return r.data
}