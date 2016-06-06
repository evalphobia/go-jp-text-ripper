package writer

func newDummyWriter() writer {
	return &dummyWriter{}
}

type dummyWriter struct{}

func (_ *dummyWriter) Write(s []string) error {
	return nil
}

func (_ *dummyWriter) Flush() {
}
