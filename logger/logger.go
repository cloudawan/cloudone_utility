// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"code.google.com/p/log4go"
	"os"
	"runtime"
)

const (
	rootPath  = "/var/log"
	logSuffix = ".log"
)

type Log struct {
	logger log4go.Logger
}

func (log Log) Critical(args ...interface{}) error {
	return log.logger.Critical(args)
}

func (log Log) Debug(args ...interface{}) {
	log.logger.Debug(args)
}

func (log Log) Error(args ...interface{}) error {
	return log.logger.Error(args)
}

func (log Log) Fine(args ...interface{}) {
	log.logger.Fine(args)
}

func (log Log) Finest(args ...interface{}) {
	log.logger.Finest(args)
}

func (log Log) Info(args ...interface{}) {
	log.logger.Info(args)
}

func (log Log) Trace(args ...interface{}) {
	log.logger.Trace(args)
}

func (log Log) Warn(args ...interface{}) error {
	return log.logger.Warn(args)
}

type LogManager struct {
	logMap        map[string]Log
	directoryPath string
}

func CreateLogManager(programName string) (*LogManager, error) {
	directoryPath, err := createDirectoryIfNotExist(programName)
	if err != nil {
		return nil, err
	}
	return &LogManager{make(map[string]Log), directoryPath}, nil
}

func createDirectoryIfNotExist(programName string) (string, error) {
	directoryPath := rootPath + string(os.PathSeparator) + programName
	return directoryPath, os.MkdirAll(directoryPath, os.ModePerm)
}

func (logManager *LogManager) GetLog(moduleName string) Log {
	filePath := logManager.directoryPath + string(os.PathSeparator) + moduleName + logSuffix
	return logManager.getLogWithFileName(filePath)
}

func (logManager *LogManager) getLogWithFileName(filePath string) Log {
	log, ok := logManager.logMap[filePath]
	if ok {
		return log
	} else {
		// Create the empty logger
		logger := make(log4go.Logger)
		fileWriter := log4go.NewFileLogWriter(filePath, false)
		fileWriter.SetFormat("[%D %T] [%L] (%S) %M")
		fileWriter.SetRotate(false)
		fileWriter.SetRotateSize(100 * 1024 * 1024)
		fileWriter.SetRotateLines(0)
		fileWriter.SetRotateDaily(false)
		logger.AddFilter("file", log4go.DEBUG, fileWriter)

		log = Log{logger}
		logManager.logMap[filePath] = log
		return log
	}
}

// Self logger
var logManager *LogManager

func init() {
	var err error
	logManager, err = CreateLogManager("kubernetes_management_utility")
	if err != nil {
		panic(err)
	}
}

func GetLog(moduleName string) Log {
	return logManager.GetLog(moduleName)
}

func GetStackTrace(maxByteAmount int, allRoutines bool) string {
	trace := make([]byte, maxByteAmount)
	count := runtime.Stack(trace, allRoutines)
	return string(trace[:count])
}
