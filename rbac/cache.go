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
	"time"
)

type Cache struct {
	Token       string
	User        *User
	CreatedTime time.Time
	ExpiredTime time.Time
}

var cacheMap map[string]*Cache = make(map[string]*Cache)

func SetCache(token string, user *User, ttl time.Duration) {
	createdTime := time.Now()
	expiredTime := createdTime.Add(ttl)

	cacheMap[token] = &Cache{
		token,
		user,
		createdTime,
		expiredTime,
	}
}

func GetCache(token string) *User {
	cache := cacheMap[token]
	if cache == nil {
		return nil
	} else {
		return cache.User
	}
}

func CheckCacheTimeout() {
	now := time.Now()
	for key, value := range cacheMap {
		if now.After(value.ExpiredTime) {
			delete(cacheMap, key)
		}
	}
}

func GetAllTokenExpiredTime() map[string]time.Time {
	expiredMap := make(map[string]time.Time)

	for key, value := range cacheMap {
		expiredMap[key] = value.ExpiredTime
	}

	return expiredMap
}
