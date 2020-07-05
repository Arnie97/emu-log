package common

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Catch(err *error) {
	if r := recover(); r != nil {
		*err = r.(error)
	}
}

func Must(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
