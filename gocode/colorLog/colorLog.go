// SEE https://github.com/liyu4/chill

package colorLog

import (
	"fmt"
	"log"
)

const (
	colorRed = uint8(iota + 91)
	colorGreen
	colorYellow
	colorBlue
	colorMagenta //洋红

	info = "[INFO]"
	trac = "[TRAC]"
	erro = "[ERRO]"
	warn = "[WARN]"
	succ = "[SUCC]"
)

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13

// Trace equals to log.Println with a yellow [TRAC] prefix
func Trace(format string, a ...interface{}) {
	prefix := yellow(trac)
	log.Println(prefix, fmt.Sprintf(format, a...))
}

// Info equals to log.Println with a blue [INFO] prefix
func Info(format string, a ...interface{}) {
	prefix := blue(info)
	log.Println(prefix, fmt.Sprintf(format, a...))
}

// Success equals to log.Println with a green [TRAC] prefix
func Success(format string, a ...interface{}) {
	prefix := green(succ)
	log.Println(prefix, fmt.Sprintf(format, a...))
}

// Warning equals to log.Println with a magenta [WARN] prefix
func Warning(format string, a ...interface{}) {
	prefix := magenta(warn)
	log.Println(prefix, fmt.Sprintf(format, a...))
}

// Error equals to log.Println with a red [ERRO] prefix
func Error(format string, a ...interface{}) {
	prefix := red(erro)
	log.Println(prefix, fmt.Sprintf(format, a...))
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorRed, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorGreen, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorYellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorBlue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorMagenta, s)
}
