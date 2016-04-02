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
	"golang.org/x/crypto/sha3"
)

type User struct {
	Name            string
	EncodedPassword string
	RoleSlice       []*Role
	ResourceSlice   []*Resource
	Description     string
}

func CreateUser(name string, password string, roleSlice []*Role, resourceSlice []*Resource, description string) *User {
	return &User{
		name,
		encodePassword(password),
		roleSlice,
		resourceSlice,
		description,
	}
}

func encodePassword(password string) string {
	fixedSlice := sha3.Sum512([]byte(password))
	byteSlice := make([]byte, 64)
	for i := 0; i < len(fixedSlice); i++ {
		byteSlice[i] = fixedSlice[i]
	}
	return hex.EncodeToString(byteSlice)
}

func (user *User) CheckPassword(password string) bool {
	return user.EncodedPassword == encodePassword(password)
}

func (user *User) HasPermission(component string, method string, path string) bool {
	for _, role := range user.RoleSlice {
		if role.HasPermission(component, method, path) {
			return true
		}
	}

	return false
}

func (user *User) HasResource(component string, path string) bool {
	for _, resource := range user.ResourceSlice {
		if resource.HasResource(component, path) {
			return true
		}
	}

	return false
}

func (user *User) CopyPartialUserDataForComponent(component string) *User {
	newUser := &User{}
	newUser.Name = user.Name
	newUser.EncodedPassword = "******"
	newUser.RoleSlice = make([]*Role, 0)
	newUser.ResourceSlice = make([]*Resource, 0)
	newUser.Description = user.Description

	for _, resource := range user.ResourceSlice {
		if resource.Component == "*" || resource.Component == component {
			newUser.ResourceSlice = append(newUser.ResourceSlice, resource)
		}
	}

	for _, role := range user.RoleSlice {
		newRole := &Role{}
		newRole.Name = role.Name
		newRole.PermissionSlice = make([]*Permission, 0)
		newRole.description = role.description

		for _, permission := range role.PermissionSlice {
			if permission.Component == "*" || permission.Component == component {
				newRole.PermissionSlice = append(newRole.PermissionSlice, permission)
			}
		}

		if len(newRole.PermissionSlice) > 0 {
			newUser.RoleSlice = append(newUser.RoleSlice, newRole)
		}
	}

	return newUser
}
