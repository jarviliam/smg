package logger

import "errors"

type Logger struct {
	title string
}

var ErrLoggerTitle = errors.New("Logger Title cannot be empty")

func NewLogger(title string) (*Logger, error) {
	if title == "" {
		return nil, ErrLoggerTitle
	}
	l := &Logger{
		title: title,
	}
	return l, nil
}

//TODO: Handle Variables
func (l *Logger) LogInfo(msg string, varadic ...interface{}) {
	e := NewLogEntry(msg)
	//TODO: Log Module here
	if e != nil {
		return
	}
}
func (l *Logger) LogError(error error) {

}
