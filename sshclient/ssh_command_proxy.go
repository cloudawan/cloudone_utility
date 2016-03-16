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
	"errors"
	"github.com/cloudawan/cloudone_utility/ioutility"
	"golang.org/x/crypto/ssh"
	"io"
	"strconv"
	"strings"
	"time"
)

func CreateSSHCommandProxy(
	dialTimeout time.Duration,
	host string,
	port int,
	user string,
	password string,
	screenHeight int,
	screenWidth int,
	interactiveMap map[string]string) *SSHCommandProxy {
	return &SSHCommandProxy{
		dialTimeout,
		host,
		port,
		user,
		password,
		screenHeight,
		screenWidth,
		interactiveMap,
		nil,
		nil,
		nil,
		nil,
		nil,
		false,
	}
}

type SSHCommandProxy struct {
	dialTimeout    time.Duration
	host           string
	port           int
	user           string
	password       string
	screenHeight   int
	screenWidth    int
	interactiveMap map[string]string
	connection     *ssh.Client
	session        *ssh.Session
	inputChannel   chan<- string
	outputChannel  <-chan string
	errorChannel   chan string
	isConnected    bool
}

func (sshCommandProxy *SSHCommandProxy) Connect() (chan<- string, <-chan string, chan string, error) {
	if sshCommandProxy.isConnected {
		return nil, nil, nil, errors.New("Already connected")
	}

	sshConfig := &ssh.ClientConfig{
		User: sshCommandProxy.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshCommandProxy.password),
		},
	}

	connection, err := dialWithTimeout("tcp", sshCommandProxy.host+":"+strconv.Itoa(sshCommandProxy.port), sshConfig, sshCommandProxy.dialTimeout)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", sshCommandProxy.screenHeight, sshCommandProxy.screenWidth, modes); err != nil {
		return nil, nil, nil, err
	}
	w, err := session.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	r, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	e, err := session.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	inputChannel, outputChannel, errorChannel := proxy(w, r, e, sshCommandProxy.interactiveMap)

	if err := session.Shell(); err != nil {
		return nil, nil, nil, err
	}

	// Ignore the ssh tty welcome page
	//<-outputChannel

	sshCommandProxy.connection = connection
	sshCommandProxy.session = session
	sshCommandProxy.inputChannel = inputChannel
	sshCommandProxy.outputChannel = outputChannel
	sshCommandProxy.errorChannel = errorChannel
	sshCommandProxy.isConnected = true

	return sshCommandProxy.inputChannel, sshCommandProxy.outputChannel, sshCommandProxy.errorChannel, nil
}

func (sshCommandProxy *SSHCommandProxy) Disconnect() error {
	if sshCommandProxy.isConnected == false {
		return errors.New("Not connected")
	} else {
		// End terminal
		sshCommandProxy.inputChannel <- "exit\n"
		// Close input
		close(sshCommandProxy.inputChannel)

		sshCommandProxy.isConnected = false
		sshCommandProxy.session.Close()
		sshCommandProxy.connection.Close()
		return nil
	}
}

func proxy(w io.Writer, r io.Reader, e io.Reader, interactiveMap map[string]string) (chan<- string, <-chan string, chan string) {
	inputChannel := make(chan string, 1)
	outputChannel := make(chan string, 1)
	errorChannel := make(chan string, 1024)

	// Issue command
	go func() {
		for cmd := range inputChannel {
			w.Write([]byte(cmd))
		}
	}()

	// Handle responsed error
	go func() {
		for {
			text, _, err := ioutility.ReadText(e, 1024*64)
			if err == io.EOF {
				close(errorChannel)
				return
			} else if err != nil {
				errorChannel <- err.Error()
				close(errorChannel)
				return
			}

			// Upon receiveing from stderr, send to error
			errorChannel <- text
		}
	}()

	// Handle responsed output
	go func() {
		for {
			text, _, err := ioutility.ReadText(r, 1024*16)
			if err == io.EOF {
				outputChannel <- text
				close(outputChannel)
				return
			} else if err != nil {
				outputChannel <- err.Error()
				close(outputChannel)
				return
			}

			for key, value := range interactiveMap {
				if strings.Contains(text, key) {
					w.Write([]byte(value))
					break
				}
			}

			outputChannel <- text
		}
	}()
	return inputChannel, outputChannel, errorChannel
}
