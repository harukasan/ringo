/*
Package debug implements utility functions to output debug message.

The debug messages can enable by Enable variable, for instance:

	debug.Enable = true

	debug.Print("debug message")
	// Output:
	// debug message

*/
package debug

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// Enable specifies whether the debug prints messages into the Output.
var Enable = false

// Output specifies the destination to print out messages.
var Output = os.Stderr

func output(s string) {
	_, file, line, _ := runtime.Caller(2)
	fmt.Fprintf(Output, "%s:%d> %s", path.Base(file), line, s)
	if s[len(s)-1] != '\n' {
		Output.Write([]byte{'\n'})
	}
}

// Print prints debug message in manner of fmt.Sprint.
func Print(v ...interface{}) {
	if Enable {
		output(fmt.Sprint(v...))
	}
}

// Println prints debug message in manner of fmt.Sprintln.
func Println(v ...interface{}) {
	if Enable {
		output(fmt.Sprintln(v...))
	}
}

// Printf prints debug message in manner of fmt.Sprintf.
func Printf(format string, v ...interface{}) {
	if Enable {
		output(fmt.Sprintf(format, v...))
	}
}
