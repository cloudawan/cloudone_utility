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

package audit

import (
	"time"
)

type AuditLog struct {
	Component         string
	Kind              string
	Path              string
	UserName          string
	CreatedTime       time.Time
	QueryParameterMap map[string][]string
	PathParameterMap  map[string]string
	RequestMethod     string
	RequestURI        string
	RequestBody       string
	RequestHeader     map[string][]string
	Description       string
}

var descriptionMap map[string]string = make(map[string]string)

func AddDescription(methodAndPath string, description string) {
	descriptionMap[methodAndPath] = description
}

func CreateAuditLog(component string, path string, userName string,
	queryParameterMap map[string][]string, pathParameterMap map[string]string,
	requestMethod string, requestURI string, requestBody string, requestHeader map[string][]string) *AuditLog {

	return &AuditLog{
		component,
		getKind(requestMethod, path),
		path,
		userName,
		time.Now(),
		queryParameterMap,
		pathParameterMap,
		requestMethod,
		requestURI,
		requestBody,
		requestHeader,
		getDescriptionFromMethodAndPath(requestMethod, path),
	}
}

func getKind(method string, path string) string {
	return method + " " + path
}

func getDescriptionFromMethodAndPath(method string, path string) string {
	return descriptionMap[getKind(method, path)]
}
