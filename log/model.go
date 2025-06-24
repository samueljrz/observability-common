package log

type Entry struct {
	Component string
	Operation string
	Message   string
	Err       error
	Fields    map[string]string

	stacktrace     string
	stacktraceHash *string
}
