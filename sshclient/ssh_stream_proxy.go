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
	"golang.org/x/crypto/ssh"
	"io"
	"strconv"
	"time"
)

func CreateSSHStreamProxy(
	dialTimeout time.Duration,
	host string,
	port int,
	user string,
	password string,
	screenHeight int,
	screenWidth int) *SSHStreamProxy {
	return &SSHStreamProxy{
		dialTimeout,
		host,
		port,
		user,
		password,
		screenHeight,
		screenWidth,
		nil,
		nil,
		false,
	}
}

type SSHStreamProxy struct {
	dialTimeout  time.Duration
	host         string
	port         int
	user         string
	password     string
	screenHeight int
	screenWidth  int
	connection   *ssh.Client
	session      *ssh.Session
	isConnected  bool
}

func (sshStreamProxy *SSHStreamProxy) Connect() (io.WriteCloser, io.Reader, io.Reader, error) {
	if sshStreamProxy.isConnected {
		return nil, nil, nil, errors.New("Already connected")
	}

	sshConfig := &ssh.ClientConfig{
		User: sshStreamProxy.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshStreamProxy.password),
		},
	}

	connection, err := dialWithTimeout("tcp", sshStreamProxy.host+":"+strconv.Itoa(sshStreamProxy.port), sshConfig, sshStreamProxy.dialTimeout)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", sshStreamProxy.screenHeight, sshStreamProxy.screenWidth, modes); err != nil {
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

	if err := session.Shell(); err != nil {
		return nil, nil, nil, err
	}

	sshStreamProxy.isConnected = true
	sshStreamProxy.session = session
	sshStreamProxy.connection = connection

	return w, r, e, err
}

func (sshStreamProxy *SSHStreamProxy) Disconnect() error {
	if sshStreamProxy.isConnected == false {
		return errors.New("Not connected")
	} else {
		sshStreamProxy.isConnected = false
		sshStreamProxy.session.Close()
		sshStreamProxy.connection.Close()
		return nil
	}
}
