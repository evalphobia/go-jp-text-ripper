package ripper

// Plugin outputs extra column with custom logic
type Plugin struct {
	Title string
	Fn    func(*TextData) string
}
