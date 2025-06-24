package config

// Mode is an enum for describing the expected behavior of the logger. The following modes are allowed:
// 	1. Noop: uses an io.Discard under the hood, so it does nothing with the logs. Useful for testing and benchmarking.
// 	2. Local: writes logs to the application's stdout. Useful for local development and debugging.
//	3. Debug: writes logs to a collector that is running locally. Useful for investigating problems as the connection with
//	the collector and parsing.
// 	3. Development: sends logs to garden's o11y using our syslog collector beta endpoint. Should be used for environments
//	in alpha/beta/dev deployments.
//	4. Production: sends logs to garden's o11y using our syslog collector endpoint. Should be used for environments in
//	prod deployments
type Mode int8

const (
	Noop Mode = iota - 1
	Local
	Debug
	Development
	Production
)
