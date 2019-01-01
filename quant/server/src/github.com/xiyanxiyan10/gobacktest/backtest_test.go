package gobacktest

import (
	"fmt"
	"testing"
)

//1 方法必须传入testing.T 2 方法名必须是以Test开头，驼峰命名
func TestBack(t *testing.T) {
	fmt.Println("create new backtest")
	backConfig := make(map[string]string)
	back := NewBackTest(backConfig)

	//@todo need scripts
	back.SetScripts("")
	//start back

	go func() {
		back.Start()
	}()

	//todo read data from history and put into back
	go func() {

	}()

	//try to print status of the backtest

}
