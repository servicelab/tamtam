/*
Copyright 2018, Eelco Cramer and the TamTam contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"github.com/clockworksoul/smudge"
	"github.com/rs/zerolog/log"
)

// SmudgeLogger allows us to use our own logger with SmudgeLogger
type SmudgeLogger struct {
}

// Log writes a log message of a certain level to the logger
func (s SmudgeLogger) Log(level smudge.LogLevel, a ...interface{}) (n int, err error) {
	str := fmt.Sprint(a...)
	switch level {
	case smudge.LogFatal:
		log.Fatal().Msg(str)
	case smudge.LogInfo:
		log.Info().Msg(str)
	case smudge.LogError:
		log.Error().Msg(str)
	case smudge.LogWarn:
		log.Warn().Msg(str)
	default:
		log.Debug().Msg(str)
	}
	return 0, nil
}

// Logf writes a log message of a certain level to the logger
func (s SmudgeLogger) Logf(level smudge.LogLevel, format string, a ...interface{}) (n int, err error) {
	switch level {
	case smudge.LogFatal:
		log.Fatal().Msgf(format, a...)
	case smudge.LogInfo:
		log.Info().Msgf(format, a...)
	case smudge.LogError:
		log.Error().Msgf(format, a...)
	case smudge.LogWarn:
		log.Warn().Msgf(format, a...)
	default:
		log.Debug().Msgf(format, a...)
	}
	return 0, nil
}
