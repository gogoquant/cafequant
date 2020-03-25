package pythonbind

import "C"
import (
	"fmt"
	"github.com/sbinet/go-python"
	"time"
)

type Go2Python struct {
	fileName string
}

func (g *Go2Python) Example(self, args *python.PyObject) *python.PyObject {
	/*
		if !python.PyString_Check(args){
			return python.Py_None
		}

		fmt.Printf("simple.example: %v", args)
	*/
	fmt.Printf("simple.example: %v", "test")
	return python.Py_None
}

// Run run the scripts
func (g *Go2Python) Run() error {
	if err := python.Initialize(); err != nil {
		return err
	}

	methods := []python.PyMethodDef{
		{"example", g.Example, python.MethNoArgs, "example function"},
	}

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
