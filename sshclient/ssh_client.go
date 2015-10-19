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

package sshclient

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

func InteractiveSSH(timeout time.Duration, host string, port int, user string, password string,
	commandSlice []string, interactiveMap map[string]string) ([]string, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	connection, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), sshConfig)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	w, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	r, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	e, err := session.StderrPipe()
	if err != nil {
		return nil, err
	}

	inputChannel, outputChannel, errorChannel := shell(w, r, e, interactiveMap)
	if err := session.Shell(); err != nil {
		return nil, err
	}

	isTimeout := false
	go func() {
		// Timeout the session to prevent got stuck
		time.Sleep(timeout)
		isTimeout = true
		session.Close()
	}()

	// Ignore the ssh tty welcome page
	<-outputChannel

	resultSlice := make([]string, 0)
	for _, command := range commandSlice {
		inputChannel <- command
		result, ok := <-outputChannel
		if ok {
			resultSlice = append(resultSlice, result)
		} else {
			break
		}
	}

	// End terminal
	inputChannel <- "exit\n"
	// Close input
	close(inputChannel)

	// Wait until I/O is closed
	session.Wait()

	buffer := bytes.Buffer{}
	for {
		errorMessage, ok := <-errorChannel
		if ok {
			buffer.WriteString(errorMessage)
		} else {
			break
		}
	}

	if isTimeout {
		buffer.WriteString("Session timeout")
	}

	if buffer.Len() > 0 {
		return resultSlice, errors.New(buffer.String())
	} else {
		return resultSlice, nil
	}
}

func shell(w io.Writer, r io.Reader, e io.Reader, interactiveMap map[string]string) (chan<- string, <-chan string, chan string) {
	inputChannel := make(chan string, 1)
	outputChannel := make(chan string, 1)
	errorChannel := make(chan string, 1024)

	waitGroup := sync.WaitGroup{}
	// Start from read
	waitGroup.Add(1)

	// Issue command
	go func() {
		for cmd := range inputChannel {
			waitGroup.Add(1)
			w.Write([]byte(cmd))
			waitGroup.Wait()
		}
	}()

	// Handle responsed error
	go func() {
		buf := make([]byte, 1024*64)
		for {
			n, err := e.Read(buf)
			if err != nil && err.Error() == "EOF" {
				close(errorChannel)
				return
			} else if err != nil {
				errorChannel <- err.Error()
				close(errorChannel)
				return
			}

			// Upon receiveing from stderr, send to error
			errorChannel <- string(buf[:n])
		}
	}()

	// Handle responsed output
	go func() {
		buf := make([]byte, 1024*64)
		length := 0
		for {
			n, err := r.Read(buf[length:])
			if err != nil && err.Error() == "EOF" {
				outputChannel <- string(buf[:length])
				close(outputChannel)
				return
			} else if err != nil {
				outputChannel <- err.Error()
				close(outputChannel)
				return
			}

			interactive := false
			currentResponse := string(buf[length:])
			for key, value := range interactiveMap {
				if strings.Contains(currentResponse, key) {
					w.Write([]byte(value))
					interactive = true
					break
				}
			}

			if interactive {
				// Ignore the response for output
			} else {
				length += n
			}

			// Keep buffing until the end of this interactive command.
			// $ is the terminal symbol where is used to tell user to enter next command.
			if length-2 > 0 && buf[length-2] == '$' {
				outputChannel <- string(buf[:length-n])
				length = 0
				waitGroup.Done()
			}
		}
	}()
	return inputChannel, outputChannel, errorChannel
}
