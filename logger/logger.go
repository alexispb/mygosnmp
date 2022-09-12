package logger

import "fmt"

type Log interface {
	Write(msg string)
	Writef(format string, args ...interface{})
}

type console struct{}

func Console() Log {
	return console{}
}

func (console) Write(msg string) {
	fmt.Println(msg)
}

func (console) Writef(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
