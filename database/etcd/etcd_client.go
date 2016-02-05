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
	"golang.org/x/net/context"
	"time"
)

var log = logger.GetLog("utility")

type EtcdClient struct {
	KeysAPI                     client.KeysAPI
	EtcdEndpoints               []string
	EtcdHeaderTimeoutPerRequest time.Duration
	EtcdBasePath                string
}

func CreateEtcdClient(etcdEndpoints []string, etcdHeaderTimeoutPerRequest time.Duration, etcdBasePath string) *EtcdClient {
	etcdClient := &EtcdClient{nil, etcdEndpoints, etcdHeaderTimeoutPerRequest, etcdBasePath}
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

func (etcdClient *EtcdClient) CreateDirectoryIfNotExist(key string) error {
	keysAPI, err := etcdClient.GetKeysAPI()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = keysAPI.Get(context.Background(), key, nil)
	errorData, ok := err.(client.Error)
	if ok == false {
		log.Error("Fail to convert error: %v", err)
		return err
	}
	// Not existing, create the directory
	if errorData.Code == client.ErrorCodeKeyNotFound {
		_, err = keysAPI.Set(context.Background(), key, "", &client.SetOptions{Dir: true})
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
