package log

import (
	"log"
)

// StdLogger use standard log package.
type StdLogger struct {
}

// Debugf logging debug information.
func (l *StdLogger) Debugf(prefix, format string, v ...interface{}) {
	log.Printf("[DEBUG]\t["+prefix+"]\t"+format, v...)
}

// Infof logging information.
func (*StdLogger) Infof(prefix, format string, v ...interface{}) {
	log.Printf("[INFO]\t["+prefix+"]\t"+format, v...)
}

// Errorf logging error information.
func (*StdLogger) Errorf(prefix, format string, v ...interface{}) {
	log.Printf("[ERROR]\t["+prefix+"]\t"+format, v...)
}
