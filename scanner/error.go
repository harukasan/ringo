package scanner

import "fmt"

// ScanError holds an error which is caused by scanner.
type ScanError struct {
	Pos int
	Err error
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("%v at pos=%d", e.Err, e.Pos)
}
