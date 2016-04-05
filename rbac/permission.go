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

type Permission struct {
	Name      string
	Component string
	Method    string
	Path      string // Path is hierarchy
}

func CreatePermission(component string, method string, path string) (*Permission, error) {
	name, err := GetPermissionName(component, method, path)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Permission{
		name,
		component,
		method,
		path,
	}, nil
}

func GetPermissionName(component string, method string, path string) (string, error) {
	if component == "" {
		log.Error("Component couldn't be empty")
		return "", errors.New("Component couldn't be empty")
	}
	if method == "" {
		log.Error("Method couldn't be empty")
		return "", errors.New("Method couldn't be empty")
	}
	if path == "" {
		log.Error("Path couldn't be empty")
		return "", errors.New("Path couldn't be empty")
	}
	if strings.Contains(path, "/") == false {
		log.Error("Path format is invalid")
		return "", errors.New("Path format is invalid")
	}

	return hex.EncodeToString([]byte(component + " " + method + " " + path)), nil
}

func (permission *Permission) HasPermission(component string, method string, path string) bool {
	// * means all
	if permission.Component == "*" {
		return true
	} else if permission.Component == component {
		// * means all
		if permission.Method == "*" {
			// * means all
			if permission.Path == "*" {
				return true
			} else {
				// Prefix for hierarchy authorization
				return strings.HasPrefix(path, permission.Path)
			}
		} else if permission.Method == method {
			// * means all
			if permission.Path == "*" {
				return true
			} else {
				// Prefix for hierarchy authorization
				return strings.HasPrefix(path, permission.Path)
			}
		} else {
			// Different method won't apply path hierarchy authorization
			return false
		}
	} else {
		return false
	}
}

// Check whether user has the target permission node or child permission node of the target permission node along the tree
func (permission *Permission) HasChildPermission(component string, method string, path string) bool {
	// * means all
	if permission.Component == "*" {
		return true
	} else if permission.Component == component {
		// * means all
		if permission.Method == "*" {
			// Prefix for hierarchy authorization. Unlike HasPermission, here it check whether target permission is the same or the child of the permissions owned
			return strings.HasPrefix(permission.Path, path)
		} else if permission.Method == method {
			// Prefix for hierarchy authorization. Unlike HasPermission, here it check whether target permission is the same or the child of the permissions owned
			return strings.HasPrefix(permission.Path, path)
		} else {
			// Different method won't apply path hierarchy authorization
			return false
		}
	} else {
		return false
	}
}
