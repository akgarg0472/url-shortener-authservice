package logger

type Config struct {
	Level          string
	LogToConsole   bool
	LogToFile      bool
	LogFilePath    string
	Enabled        bool
	IsInfoEnabled  bool
	IsErrorEnabled bool
	IsDebugEnabled bool
	IsTraceEnabled bool
	IsFatalEnabled bool
}
