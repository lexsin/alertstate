// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package alertstate

import (
	"fmt"
	//"ifm/modules/mconfig"
	"os"
	"path"

	//"path/filepath"
	//"preprocess/modules/mlog"
	"strings"

	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
const (
	LevelEmerg = iota
	LevelAlert
	LevelCrit
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
)

// BeeLogger references the used application logger.
var BeeLogger = logs.NewLogger(100)

// SetLevel sets the global log level used by the simple logger.
func SetLogLevel(l int) {
	BeeLogger.SetLevel(l)
}

// SetLogFuncCall set the CallDepth, default is 3
func SetLogFuncCall(b bool) {
	BeeLogger.EnableFuncCallDepth(b)
	BeeLogger.SetLogFuncCallDepth(3)
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := BeeLogger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

// Emergency logs a message at emergency level.
func Emergency(v ...interface{}) {
	BeeLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Emergf(format string, v ...interface{}) {
	BeeLogger.Emergency(format, v...)
}

// Alert logs a message at alert level.
func Alert(v ...interface{}) {
	BeeLogger.Alert(generateFmtStr(len(v)), v...)
}

func Alertf(format string, v ...interface{}) {
	BeeLogger.Alert(format, v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	BeeLogger.Critical(generateFmtStr(len(v)), v...)
}

func Critf(format string, v ...interface{}) {
	BeeLogger.Critical(format, v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	BeeLogger.Error(generateFmtStr(len(v)), v...)
}

func Errorf(format string, v ...interface{}) {
	BeeLogger.Error(format, v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	BeeLogger.Warning(generateFmtStr(len(v)), v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	BeeLogger.Warn(generateFmtStr(len(v)), v...)
}

func Warnf(format string, v ...interface{}) {
	BeeLogger.Warn(format, v...)
}

// Notice logs a message at notice level.
func Notice(v ...interface{}) {
	BeeLogger.Notice(generateFmtStr(len(v)), v...)
}

func Noticef(format string, v ...interface{}) {
	BeeLogger.Notice(format, v...)
}

// Informational logs a message at info level.
func Informational(v ...interface{}) {
	BeeLogger.Informational(generateFmtStr(len(v)), v...)
}

// Info compatibility alias for Warning()
func Info(v ...interface{}) {
	BeeLogger.Info(generateFmtStr(len(v)), v...)
}

func Infof(format string, v ...interface{}) {
	BeeLogger.Info(format, v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	BeeLogger.Debug(generateFmtStr(len(v)), v...)
}

func Debugf(format string, v ...interface{}) {
	BeeLogger.Debug(format, v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	BeeLogger.Trace(generateFmtStr(len(v)), v...)
}

func Tracef(format string, v ...interface{}) {
	BeeLogger.Trace(format, v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func GetMlogLevel(input string) int {
	level := 0
	switch strings.ToLower(input) {
	case "emerg":
		level = logs.LevelEmergency
	case "crit":
		level = logs.LevelCritical
	case "error":
		level = logs.LevelError
	case "warning":
		level = logs.LevelWarning
	case "notice":
		level = logs.LevelNotice
	case "info":
		level = logs.LevelInformational
	case "debug":
		level = logs.LevelDebug
	default:
		level = logs.LevelInformational
	}
	return level
}

type LevelString string

func (me LevelString) TransLevel() int {
	switch string(me) {
	case "emerg":
		return LevelEmerg
	case "alert":
		return LevelAlert
	case "crit":
		return LevelCrit
	case "error":
		return LevelError
	case "warning":
		return LevelWarning
	case "notice":
		return LevelNotice
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		return LevelInfo
	}
}

func CreateDir(dir string) error {
	if absolute := path.IsAbs(dir); !absolute {
		return nil
	}
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	} else {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0777); err != nil {
				return err
			}
		}
	}
	return err
}

func initLog() {
	dir := "logs/alertstate"
	CreateDir(dir)
	if err := SetLogger("file", `{"filename":"logs/alertstate/server.log"}`); err != nil {
		fmt.Println(err)
	}
	//SetLogger("console", "")
	//temp, _ := mconfig.Conf.String("log", "level")
	//SetLogLevel(LevelString(temp).TransLevel())
	SetLogLevel(LevelDebug)
	BeeLogger.Async()
}
