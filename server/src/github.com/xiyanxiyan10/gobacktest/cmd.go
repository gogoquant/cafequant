package gobacktest

type Cmd struct {
	Event
	cmd string
}

func (c *Cmd) SetCmd(cmd string) {
	c.cmd = cmd
}

func (c *Cmd) Cmd() string {
	return c.cmd
}
