package pythonbind

import "C"
import (
	"github.com/sbinet/go-python"
	"time"
)

type Go2Python struct {
	fileName string
}

func Hello() int {
	return 3
}

// Run run the scripts
func (g Go2Python) Run() error {
	if err := python.Initialize(); err != nil {
		return err
	}
	//python.py()
	methods := make([]python.PyMethodDef, 1)
	method := &methods[0]
	//method.Meth
	method.Name = "hello"
	method.Doc = "hello world"
	method.Flags = python.MethNoArgs

	_, err := python.Py_InitModule("snack", methods)
	if err != nil {
		panic(err)
	}
	if err := python.PyRun_SimpleFile(g.fileName); err != nil {
		return err
	}
	if err := python.Finalize(); err != nil {
		return err
	}
	time.Sleep(time.Minute * 5)
	return nil
}
