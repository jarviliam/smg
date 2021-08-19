package writer

type Writer struct {
	MaxLines int
}

var (
	defaultMaxLines = 100
)

func NewWriter() *Writer {
	w := &Writer{
		MaxLines: defaultMaxLines,
	}
	return w
}
