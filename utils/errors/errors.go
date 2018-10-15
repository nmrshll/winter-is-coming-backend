package errors

import (
	"fmt"

	"github.com/fatih/color"
)

func Log(err error) {
	fmt.Println(color.RedString("[ERROR] %v\n", err))
}

func Println(format string, values ...interface{}) {
	fmt.Printf(color.RedString("[ERROR] %s\n", fmt.Sprintf(format, values...)))
}

func Wrap(err error, format string, args ...interface{}) error {
	return fmt.Errorf(format+":\n"+err.Error(), args...)
}
