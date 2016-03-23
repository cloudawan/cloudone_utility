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

package deepcopy

import ()

func DeepCopyJsonData(dataFrom interface{}) interface{} {
	switch dataFrom.(type) {
	case map[string]interface{}:
		dataTo := make(map[string]interface{})
		DeepOverwriteJsonMap(dataFrom.(map[string]interface{}), dataTo)
		return dataTo
	case []interface{}:
		fromJsonSlice := dataFrom.([]interface{})
		toJsonSlice := make([]interface{}, 0)
		DeepOverwriteJsonSlice(&fromJsonSlice, &toJsonSlice)
		return toJsonSlice
	default:
		return dataFrom
	}
}

func DeepOverwriteJsonMap(mapFrom map[string]interface{}, mapTo map[string]interface{}) {
	for key, value := range mapFrom {
		switch value.(type) {
		case map[string]interface{}:
			switch mapTo[key].(type) {
			case map[string]interface{}:
			default:
				// If not json map, abandom the old value and create an empty json map
				mapTo[key] = make(map[string]interface{})
			}
			DeepOverwriteJsonMap(value.(map[string]interface{}), mapTo[key].(map[string]interface{}))
		case []interface{}:
			switch mapTo[key].(type) {
			case []interface{}:
			default:
				// If not json slice, abandom the old value and create an empty json slice
				mapTo[key] = make([]interface{}, 0)
			}
			fromJsonSlice := value.([]interface{})
			toJsonSlice := mapTo[key].([]interface{})
			DeepOverwriteJsonSlice(&fromJsonSlice, &toJsonSlice)
			// Reassign again because the slice may have append operation which cause the pointer address to change
			mapTo[key] = toJsonSlice
		default:
			mapTo[key] = value
		}
	}
}

func DeepOverwriteJsonSlice(sliceFromPointer *[]interface{}, sliceToPointer *[]interface{}) {
	sliceFrom := *sliceFromPointer
	for index, value := range sliceFrom {
		sliceTo := *sliceToPointer
		sliceToLength := len(sliceTo)
		switch value.(type) {
		case map[string]interface{}:
			if index < sliceToLength {
				// Index exist
				switch sliceTo[index].(type) {
				case map[string]interface{}:
				default:
					// Not map so create a new one
					sliceTo[index] = make(map[string]interface{})
				}
				DeepOverwriteJsonMap(value.(map[string]interface{}), sliceTo[index].(map[string]interface{}))
			} else {
				// Index not exist
				newMap := make(map[string]interface{})
				DeepOverwriteJsonMap(value.(map[string]interface{}), newMap)
				*sliceToPointer = append(sliceTo, newMap)
			}
		case []interface{}:
			if index < sliceToLength {
				// Index exist
				switch sliceTo[index].(type) {
				case []interface{}:
				default:
					// Not slice so create a new one
					sliceTo[index] = make([]interface{}, 0)
				}
				fromJsonSlice := value.([]interface{})
				toJsonSlice := sliceTo[index].([]interface{})
				DeepOverwriteJsonSlice(&fromJsonSlice, &toJsonSlice)
				// Reassign again because the slice may have append operation which cause the pointer address to change
				sliceTo[index] = toJsonSlice
			} else {
				// Index not exist
				newSlice := make([]interface{}, 0)
				fromJsonSlice := value.([]interface{})
				toJsonSlice := newSlice
				DeepOverwriteJsonSlice(&fromJsonSlice, &toJsonSlice)
				*sliceToPointer = append(sliceTo, newSlice)
			}
		default:
			// Not map or slice
			if index < sliceToLength {
				// Index exist, replace the old value
				sliceTo[index] = value
			} else {
				// Index not exist, create a slot if not existing
				*sliceToPointer = append(sliceTo, value)
			}
		}
	}
}
