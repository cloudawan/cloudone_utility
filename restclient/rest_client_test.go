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

package restclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestGet(t *testing.T) {
	jsonMap, _ := RequestGet("http://172.16.0.113:8080/api/v1beta3/namespaces/default/replicationcontrollers/cassandra/", true)

	fmt.Println(jsonMap.(map[string]interface{})["spec"].(map[string]interface{})["replicas"].(json.Number).Int64())
}

/*
func TestRequestGet(t *testing.T) {
	_, err := RequestGet("http://192.168.0.33:8080/api/v1beta3/namespaces/default/replicationcontrollers/nginx/")
	if err == nil {
		t.Errorf("This should be error")
	}

	a, err := RequestGet("http://192.168.0.33:8080/api/v1beta3/namespaces/default/replicationcontrollers/flask/")
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	fmt.Println(a)
}
*/
/*
type ReplicationControllerMetricList struct {
	ErrorSlice []error
	ReplicationControllerMetricSlice []ReplicationControllerMetric
}

type ReplicationControllerMetric struct {
	Namespace string
	ReplicationControllerName string
	ValidPodSlice []bool
	PodMetricSlice []PodMetric
	Size int
}

type PodMetric struct {
	KubeletHost string
	Namespace string
	PodName string
	ValidContainerSlice []bool
	ContainerMetricSlice []ContainerMetric
}

type ContainerMetric struct {
	ContainerName string
}

func TestRequestGetWithStructure(t *testing.T) {
	replicationControllerMetric, err := RequestGetWithStructure("http://127.0.0.1:8081/replicationcontrollermetric/default?kubeapihost=192.168.0.33&kubeapiport=8080", &ReplicationControllerMetricList{})

	//replicationControllerMetric, err := RequestGet("http://127.0.0.1:8081/replicationcontrollermetric/default?kubeapihost=192.168.0.33&kubeapiport=8080")
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	fmt.Println(replicationControllerMetric, err)
}
*/
/*
type Service struct {
	Name          string
	Namespace     string
	PortSlice     []ServicePort
	Selector      string
	PortalIP      string
	PublicIPSlice []string
	LabelMap      map[string]string
}

type ServicePort struct {
	Name       string
	Protocol   string
	Port       string
	TargetPort string
}

func TestRequestGetWithStructure(t *testing.T) {
	serviceSlice := make([]Service, 0)
	returnedServiceSlice, err := RequestGetWithStructure("http://192.168.0.15:8081/services/default?kubeapihost=192.168.0.33&kubeapiport=8080", serviceSlice)

	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	fmt.Println(returnedServiceSlice, err)
}
*/
