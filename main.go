package qor

import (
	"fmt"

	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

func filenameWithLineNum() string {
	var total = 10
	var results []string
	for i := 2; i < 15; i++ {
		if _, file, line, ok := runtime.Caller(i); ok {
			total--
			results = append(results[:0],
				append(
					[]string{fmt.Sprintf("%v:%v", strings.TrimPrefix(file, os.Getenv("GOPATH")+"src/"), line)},
					results[0:len(results)]...)...)

			if total == 0 {
				return strings.Join(results, "\n")
			}
		}
	}
	return ""
}

func ExitWithMsg(str string, value ...interface{}) {
	debug.PrintStack()
	fmt.Printf("\n"+filenameWithLineNum()+"\n"+str+"\n", value...)
	os.Exit(1)
}
