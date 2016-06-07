package writer

func newDummyWriter() writer {
	return &dummyWriter{}
}

// dummy writer do nothing!
type dummyWriter struct{}

// Write is dummy method
func (_ *dummyWriter) Write(s []string) error {
	return nil
}

// Flush is dummy method
func (_ *dummyWriter) Flush() {
}
