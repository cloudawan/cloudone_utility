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

/*
import (
	"fmt"
	"testing"
	"time"
)

func TestSSHProxy(t *testing.T) {
	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = "cloud4win\n"
	sshCommandProxy := CreateSSHCommandProxy(1*time.Second, 10*time.Minute, "192.168.0.31", 22, "cloudawan", "cloud4win", interactiveMap)

	err := sshCommandProxy.Connect()
	fmt.Println(err)

	inputChannel, outputChannel, _, err := sshCommandProxy.GetChannels()
	fmt.Println(err)

	inputChannel <- "sudo pwd\n"
	result, ok := <-outputChannel
	fmt.Println(ok, result)

	inputChannel <- "sudo pwd\n"
	result, ok = <-outputChannel
	fmt.Println(ok, result)

	inputChannel <- "pwd\n"
	result, ok = <-outputChannel
	fmt.Println(ok, result)

	//inputChannel <- "sudo gluster volume list\n"
	//result, ok = <-outputChannel
	//fmt.Println(ok, result)

	sshCommandProxy.Disconnect()
}
*/
