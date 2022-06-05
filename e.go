package L

import "os"

// EPanic is a Config.E that panics when there is an error.
func EPanic(_ Logger, _ *Config, e error) {
	panic(e)
}

// EFatal is a Config.E that calls ELog and then exits.
func EFatal(l Logger, c *Config, e error) {
	ELog(l, c, e)
	os.Exit(7)
}

// ELog is a Config.E that safely logs the error 'e' in a dict with key '"LE"'.
func ELog(l Logger, _ *Config, e error) {
	ev := l.Dict()
	ev.Field("LE", e.Error())
	l.Log(ev)
}
