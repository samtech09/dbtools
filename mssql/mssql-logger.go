package mssql

import "github.com/rs/zerolog"

type mssqlLogger struct {
	zerolog.Logger
}

// NewLogger accepts a zerolog.Logger as input and returns a new custom pgx
// logging fascade as output.
func newLogger(logger zerolog.Logger, connName string) *mssqlLogger {
	m := mssqlLogger{}
	m.Logger = logger.With().Str("module", connName).Logger()
	return &m
}

func (l mssqlLogger) Println(v ...interface{}) {
	l.Print(v...)
	l.Print("\n")
}

// func newLogger(isErr, isInfo, isDebug bool, callerInfo bool) *mssqlLogger {
// 	m := mssqlLogger{}
// 	m.Logger = New(isErr, isInfo, isDebug, callerInfo)
// 	return &m
// }
