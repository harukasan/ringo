package debug_test

import (
	"os"

	"github.com/harukasan/ringo/debug"
)

func ExamplePrint() {
	debug.Enable = true
	// set output to stdout for testing
	debug.Output = os.Stdout
	//
	//
	debug.Print("Hello, world!") // line: 15
	// Output:
	// debug_test.go:15> Hello, world!
}

func ExamplePrintln() {
	debug.Enable = true
	// set output to stdout for testing
	debug.Output = os.Stdout

	debug.Println("Hello, world!") // line: 25
	// Output:
	// debug_test.go:25> Hello, world!
}

func ExamplePrintf() {
	debug.Enable = true
	// set output to stdout for testing
	debug.Output = os.Stdout

	debug.Printf("a: %v", true) // line: 35
	// Output:
	// debug_test.go:35> a: true
}
