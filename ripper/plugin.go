package ripper

// Plugin outputs extra column with custom logic
type Plugin struct {
	Title string
	Fn    func(*TextData) string
}

// PostFilter outputs extra column with custom logic after plugin process
type PostFilter struct {
	Title string
	// Fn arguments is each row data
	Fn func(map[string]string) string
}

// PreFilter normalizes text data before text processing
type PreFilter struct {
	Title string
	Fn    func(string) string
}
