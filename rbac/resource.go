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

package rbac

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Resource struct {
	Name      string
	Component string
	Path      string // Path is hierarchy
}

func CreateResource(component string, path string) (*Resource, error) {
	name, err := GetResourceName(component, path)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &Resource{
		name,
		component,
		path,
	}, nil
}

func GetResourceName(component string, path string) (string, error) {
	if component == "" {
		log.Error("Component couldn't be empty")
		return "", errors.New("Component couldn't be empty")
	}
	if path == "" {
		log.Error("Path couldn't be empty")
		return "", errors.New("Path couldn't be empty")
	}

	return hex.EncodeToString([]byte(component + " " + path)), nil
}

func (resource *Resource) HasResource(component string, path string) bool {
	// * means all
	if resource.Component == "*" {
		return true
	} else if resource.Component == component {
		// * means all
		if path == "*" {
			return true
			// Prefix for hierarchy authorization
		} else if strings.HasPrefix(resource.Path, path) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
