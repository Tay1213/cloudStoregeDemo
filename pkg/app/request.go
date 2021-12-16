package app

import (
	"fmt"
	"github.com/astaxie/beego/validation"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		fmt.Printf("errKey: %#v, errMsg:%#v", err.Key, err.Message)
	}

	return
}
