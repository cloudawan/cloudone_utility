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
	"fmt"
	"testing"
	"time"
)

func TestRequestGet(t *testing.T) {
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, "sudo pwd\n")
	commandSlice = append(commandSlice, "sudo gluster --mode=script volume info\n")
	// commandSlice = append(commandSlice, "sudo gluster --mode=script volume delete test")
	// commandSlice = append(commandSlice, "sudo gluster volume create test replica 2 192.168.0.25:/data/glusterfs/test 192.168.0.26:/data/glusterfs/test force\n")
	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = "cloud4win\n"
	resultSlice, err := InteractiveSSH(3*time.Second, "192.168.0.25", 22, "cloudawan", "cloud4win", commandSlice, interactiveMap)
	fmt.Println(resultSlice)
	fmt.Println(err)
}
