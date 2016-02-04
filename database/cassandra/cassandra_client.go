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

package cassandra

import (
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/gocql/gocql"
	"time"
)

var log = logger.GetLog("utility")

type CassandraClient struct {
	session             *gocql.Session
	clusterIp           []string
	clusterPort         int
	keyspace            string
	replicationStrategy string
	timeout             time.Duration
}

func CreateCassandraClient(clusterIp []string, clusterPort int,
	keyspace string, replicationStrategy string, timeout time.Duration) *CassandraClient {
	cassandraClient := &CassandraClient{nil, clusterIp, clusterPort, keyspace, replicationStrategy, timeout}
	cassandraClient.GetSession()
	return cassandraClient
}

func (cassandraClient *CassandraClient) GetSession() (returnedSession *gocql.Session, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("GetSession Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			cassandraClient.session = nil
			returnedSession = nil
			returnedError = err.(error)
		}
	}()

	if cassandraClient.session != nil {
		return cassandraClient.session, nil
	} else {
		cluster := gocql.NewCluster(cassandraClient.clusterIp...)
		cluster.Timeout = cassandraClient.timeout
		cluster.Port = cassandraClient.clusterPort
		session, err := cluster.CreateSession()
		if err != nil {
			log.Critical("Fail to create Cassandra session: %s", err)
			cassandraClient.session = nil
			return nil, err
		} else {
			if err := session.Query("CREATE KEYSPACE IF NOT EXISTS " + cassandraClient.keyspace + " WITH replication = " + cassandraClient.replicationStrategy).Exec(); err != nil {
				log.Critical("Fail to check if not exist then create keyspace error: %s", err)
				cassandraClient.session = nil
				return nil, err
			} else {
				session.Close()
				cluster.Keyspace = cassandraClient.keyspace
				session, err := cluster.CreateSession()
				if err != nil {
					log.Critical("Fail to create Cassandra session: %s", err)
					cassandraClient.session = nil
					return nil, err
				} else {
					cassandraClient.session = session
					return cassandraClient.session, nil
				}
			}
		}
	}
}

func (cassandraClient *CassandraClient) CloseSession() {
	cassandraClient.session.Close()
	cassandraClient.session = nil
}

func (cassandraClient *CassandraClient) ResetSession() {
	cassandraClient.CloseSession()
	cassandraClient.GetSession()
}

func (cassandraClient *CassandraClient) CreateTableIfNotExist(tableSchema string, retryAmount int, retryInterval time.Duration) (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("CreateTableIfNotExist Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	session, err := cassandraClient.GetSession()
	if err != nil {
		return err
	}

	for i := 0; i < retryAmount; i++ {
		if err := session.Query(tableSchema).Exec(); err == nil {
			return nil
		} else {
			log.Error("Check if not exist then create table schema %s error: %s", tableSchema, err)
			returnedError = err
		}
		time.Sleep(retryInterval)
	}

	return returnedError
}
