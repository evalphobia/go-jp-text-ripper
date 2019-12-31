package log

// DefaultLogger is default Logger.
var DefaultLogger Logger

// Logger is logging interface.
type Logger interface {
	Debugf(prefix, format string, v ...interface{})
	Infof(prefix, format string, v ...interface{})
	Errorf(prefix, format string, v ...interface{})
}

func init() {
	v := &DummyLogger{}
	DefaultLogger = v
}

// DummyLogger does not ouput anything
type DummyLogger struct{}

// Debugf does nothing.
func (*DummyLogger) Debugf(prefix, format string, v ...interface{}) {}

// Infof does nothing.
func (*DummyLogger) Infof(prefix, format string, v ...interface{}) {}

// Errorf does nothing.
func (*DummyLogger) Errorf(prefix, format string, v ...interface{}) {}
