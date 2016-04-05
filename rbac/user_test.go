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
	"fmt"
	"testing"
)

func TestHasPermission(t *testing.T) {
	permissionSlice := make([]*Permission, 0)
	permission := &Permission{"P1", "cloudone_gui", "GET", "/gui/inventory/service"}
	permissionSlice = append(permissionSlice, permission)
	roleSlice := make([]*Role, 0)
	role := &Role{"R1", permissionSlice, ""}
	roleSlice = append(roleSlice, role)
	user := &User{"u", "p", roleSlice, nil, ""}

	fmt.Println(user.HasPermission("cloudone_gui", "GET", "/gui/inventory/service"))
	fmt.Println(user.HasPermission("cloudone_gui", "GET", "/gui/inventory"))
	fmt.Println(user.HasChildPermission("cloudone_gui", "GET", "/gui/inventory/service"))
	fmt.Println(user.HasChildPermission("cloudone_gui", "GET", "/gui/inventory"))
}
