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

package etcd

import (
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/coreos/etcd/client"
	"time"
)

var log = logger.GetLog("utility")

type EtcdClient struct {
	KeysAPI                     client.KeysAPI
	EtcdEndpoints               []string
	EtcdHeaderTimeoutPerRequest time.Duration
}

func CreateEtcdClient(etcdEndpoints []string, etcdHeaderTimeoutPerRequest time.Duration) *EtcdClient {
	etcdClient := &EtcdClient{nil, etcdEndpoints, etcdHeaderTimeoutPerRequest}
	etcdClient.GetKeysAPI()
	return etcdClient
}

func (etcdClient *EtcdClient) GetKeysAPI() (returnedKeysAPI client.KeysAPI, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("GetKeysAPI Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			etcdClient.KeysAPI = nil
			returnedKeysAPI = nil
			returnedError = err.(error)
		}
	}()

	if etcdClient.KeysAPI != nil {
		return etcdClient.KeysAPI, nil
	} else {
		config := client.Config{
			Endpoints:               etcdClient.EtcdEndpoints,
			Transport:               client.DefaultTransport,
			HeaderTimeoutPerRequest: etcdClient.EtcdHeaderTimeoutPerRequest,
		}
		configuredClient, err := client.New(config)
		if err != nil {
			etcdClient.KeysAPI = nil
			log.Error(err)
			return nil, err
		}
		keysAPI := client.NewKeysAPI(configuredClient)

		etcdClient.KeysAPI = keysAPI
		return keysAPI, nil
	}
}
