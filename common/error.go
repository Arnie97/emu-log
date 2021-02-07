package common

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// init enables a human friendly log format.
func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Catch captures a possible panic and return it as an error.
func Catch(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("panic: %v", r)
	}
}

// Must prints the error message and exit immediately if error is not nil.
func Must(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
