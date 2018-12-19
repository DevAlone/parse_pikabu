package logger

import (
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("parse_pikabu/core")
var ParserLog = logging.MustGetLogger("parse_pikabu/parser")
var LogFormat = logging.MustStringFormatter(
	`%{color}%{module} %{pid} %{level:.5s} %{time:15:04:05.000} %{shortfile} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)
