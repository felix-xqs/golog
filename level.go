package golog

type Level int8

const (
	DebugLevel Level = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var levelNames = [...]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	PanicLevel: "panic",
	FatalLevel: "fatal",
}

func (l Level) ToLowerString() string {
	return levelNames[l]
}

func LevelStringToCode(levelString string) Level {
	for i, ls := range levelNames {
		if ls == levelString {
			return Level(i)
		}
	}
	return DebugLevel
}
