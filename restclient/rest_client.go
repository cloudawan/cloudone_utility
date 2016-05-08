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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var insecureHTTPSClient *http.Client = nil

func GetInsecureHTTPSClient() *http.Client {
	// Skip the server side certificate checking
	if insecureHTTPSClient == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		insecureHTTPSClient = &http.Client{
			Transport: transport,
		}
		return insecureHTTPSClient
	} else {
		return insecureHTTPSClient
	}
}

func HealthCheck(url string, timeout time.Duration) (returnedResult bool, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedResult = false
			returnedError = err.(error)
		}
	}()
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	insecureHTTPSClient := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	response, err := insecureHTTPSClient.Do(request)
	if err != nil {
		return false, err
	} else {
		defer response.Body.Close()
		if response.StatusCode == 200 || response.StatusCode == 204 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

type RequestError struct {
	Url           string
	StatusCode    int
	ResponseData  interface{}
	ResponseBody  *string
	ReturnedError error
}

func (requestError RequestError) AllData() string {
	jsonMap := make(map[string]interface{})
	jsonMap["Url"] = requestError.Url
	jsonMap["StatusCode"] = requestError.StatusCode
	jsonMap["ResponseData"] = requestError.ResponseData
	if requestError.ResponseBody != nil {
		jsonMap["ResponseBody"] = *requestError.ResponseBody
	} else {
		jsonMap["ResponseBody"] = ""
	}
	if requestError.ReturnedError != nil {
		jsonMap["ReturnedError"] = requestError.ReturnedError.Error()
	} else {
		jsonMap["ReturnedError"] = nil
	}
	byteSlice, _ := json.Marshal(jsonMap)
	return string(byteSlice)
}

func (requestError RequestError) Error() string {
	if requestError.ResponseData != nil {
		byteSlice, err := json.Marshal(requestError.ResponseData)
		if err != nil {
			if requestError.ResponseBody != nil {
				return *requestError.ResponseBody
			} else {
				return ""
			}
		} else {
			return string(byteSlice)
		}
	} else {
		if requestError.ResponseBody != nil {
			return *requestError.ResponseBody
		} else {
			return ""
		}
	}
}

func Request(method string, url string, body interface{}, headerMap map[string]string, useJsonNumberInsteadFloat64ForResultJson bool) (returnedStatusCode int, returnedJsonMapOrJsonSlice interface{}, returnedResponseBody *string, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedStatusCode = 500
			returnedJsonMapOrJsonSlice = nil
			returnedResponseBody = nil
			returnedError = err.(error)
		}
	}()

	byteSlice, err := json.Marshal(body)
	if err != nil {
		return 500, nil, nil, err
	} else {
		var request *http.Request
		if body == nil {
			request, err = http.NewRequest(method, url, nil)
		} else {
			request, err = http.NewRequest(method, url, bytes.NewReader(byteSlice))
			request.Header.Add("Content-Type", "application/json")
		}

		for key, value := range headerMap {
			request.Header.Add(key, value)
		}

		if err != nil {
			return 500, nil, nil, err
		} else {
			response, err := GetInsecureHTTPSClient().Do(request)
			if err != nil {
				return 500, nil, nil, err
			} else {
				defer response.Body.Close()
				if response.ContentLength == 0 {
					return response.StatusCode, nil, nil, nil
				}
				responseBody, err := ioutil.ReadAll(response.Body)
				responseBodyText := string(responseBody)
				if err != nil {
					return 500, nil, &responseBodyText, err
				} else {
					var jsonMapOrJsonSlice interface{}
					if useJsonNumberInsteadFloat64ForResultJson {
						decoder := json.NewDecoder(bytes.NewReader(responseBody))
						decoder.UseNumber()
						err := decoder.Decode(&jsonMapOrJsonSlice)
						if err != nil {
							return response.StatusCode, nil, &responseBodyText, err
						} else {
							return response.StatusCode, jsonMapOrJsonSlice, &responseBodyText, nil
						}
					} else {
						err := json.Unmarshal(responseBody, &jsonMapOrJsonSlice)
						if err != nil {
							return response.StatusCode, nil, &responseBodyText, err
						} else {
							return response.StatusCode, jsonMapOrJsonSlice, &responseBodyText, nil
						}
					}
				}
			}
		}
	}
}

func RequestGet(url string, headerMap map[string]string, useJsonNumberInsteadFloat64ForResultJson bool) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := Request("GET", url, nil, headerMap, useJsonNumberInsteadFloat64ForResultJson)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestPost(url string, body interface{}, headerMap map[string]string, useJsonNumberInsteadFloat64ForResultJson bool) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := Request("POST", url, body, headerMap, useJsonNumberInsteadFloat64ForResultJson)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestPut(url string, body interface{}, headerMap map[string]string, useJsonNumberInsteadFloat64ForResultJson bool) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := Request("PUT", url, body, headerMap, useJsonNumberInsteadFloat64ForResultJson)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestDelete(url string, body interface{}, headerMap map[string]string, useJsonNumberInsteadFloat64ForResultJson bool) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := Request("DELETE", url, body, headerMap, useJsonNumberInsteadFloat64ForResultJson)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestWithStructure(method string, url string, body interface{}, returnedStructure interface{}, headerMap map[string]string) (returnedStatusCode int, returnedJsonMapOrJsonSlice interface{}, returnedResponseBody *string, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedStatusCode = 500
			returnedJsonMapOrJsonSlice = nil
			returnedError = err.(error)
			returnedResponseBody = nil
		}
	}()

	byteSlice, err := json.Marshal(body)
	if err != nil {
		return 500, nil, nil, err
	} else {
		request, err := http.NewRequest(method, url, bytes.NewReader(byteSlice))
		request.Header.Add("Content-Type", "application/json")

		for key, value := range headerMap {
			request.Header.Add(key, value)
		}

		if err != nil {
			return 500, nil, nil, err
		} else {
			response, err := GetInsecureHTTPSClient().Do(request)
			if err != nil {
				return 500, nil, nil, err
			} else {
				defer response.Body.Close()
				if response.ContentLength == 0 {
					return response.StatusCode, nil, nil, nil
				}
				responseBody, err := ioutil.ReadAll(response.Body)
				responseBodyText := string(responseBody)
				if err != nil {
					return 500, returnedJsonMapOrJsonSlice, nil, err
				} else {
					var jsonMapOrJsonSlice interface{}
					err := json.Unmarshal(responseBody, &jsonMapOrJsonSlice)
					if err != nil {
						jsonMapOrJsonSlice = nil
					}

					err = json.Unmarshal(responseBody, &returnedStructure)
					if err != nil {
						return response.StatusCode, jsonMapOrJsonSlice, &responseBodyText, err
					} else {
						return response.StatusCode, jsonMapOrJsonSlice, &responseBodyText, nil
					}
				}
			}
		}
	}
}

func RequestGetWithStructure(url string, returnedStrucutre interface{}, headerMap map[string]string) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := RequestWithStructure("GET", url, nil, returnedStrucutre, headerMap)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestPostWithStructure(url string, body interface{}, returnedStrucutre interface{}, headerMap map[string]string) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := RequestWithStructure("POST", url, body, returnedStrucutre, headerMap)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestPutWithStructure(url string, body interface{}, returnedStrucutre interface{}, headerMap map[string]string) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := RequestWithStructure("PUT", url, body, returnedStrucutre, headerMap)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestDeleteWithStructure(url string, body interface{}, returnedStrucutre interface{}, headerMap map[string]string) (interface{}, error) {
	statusCode, jsonMapOrJsonSlice, responseBody, err := RequestWithStructure("DELETE", url, body, returnedStrucutre, headerMap)
	if err != nil {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return jsonMapOrJsonSlice, nil
	} else {
		return jsonMapOrJsonSlice, RequestError{url, statusCode, jsonMapOrJsonSlice, responseBody, nil}
	}
}

func RequestByteSliceResult(method string, url string, body map[string]interface{}, headerMap map[string]string) (returnedStatusCode int, returnedByteSlice []byte, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			returnedStatusCode = 500
			returnedByteSlice = nil
			returnedError = err.(error)
		}
	}()

	byteSlice, err := json.Marshal(body)
	if err != nil {
		return 500, nil, err
	} else {
		request, err := http.NewRequest(method, url, bytes.NewReader(byteSlice))
		request.Header.Add("Content-Type", "application/json")

		for key, value := range headerMap {
			request.Header.Add(key, value)
		}

		if err != nil {
			return 500, nil, err
		} else {
			response, err := GetInsecureHTTPSClient().Do(request)
			if err != nil {
				return 500, nil, err
			} else {
				defer response.Body.Close()
				if response.ContentLength == 0 {
					return response.StatusCode, nil, nil
				}
				responseBody, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return 500, nil, err
				} else {
					return response.StatusCode, responseBody, nil
				}
			}
		}
	}
}

func RequestGetByteSliceResult(url string, headerMap map[string]string) ([]byte, error) {
	statusCode, byteSlice, err := RequestByteSliceResult("GET", url, nil, headerMap)
	text := string(byteSlice)
	if err != nil {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, err}
	} else if statusCode == 200 || statusCode == 204 {
		return byteSlice, nil
	} else {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, nil}
	}
}

func RequestPostByteSliceResult(url string, body map[string]interface{}, headerMap map[string]string) ([]byte, error) {
	statusCode, byteSlice, err := RequestByteSliceResult("POST", url, body, headerMap)
	text := string(byteSlice)
	if err != nil {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, err}
	} else if statusCode == 200 || statusCode == 201 || statusCode == 202 {
		return byteSlice, nil
	} else {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, nil}
	}
}

func RequestPutByteSliceResult(url string, body map[string]interface{}, headerMap map[string]string) ([]byte, error) {
	statusCode, byteSlice, err := RequestByteSliceResult("PUT", url, body, headerMap)
	text := string(byteSlice)
	if err != nil {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return byteSlice, nil
	} else {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, nil}
	}
}

func RequestDeleteByteSliceResult(url string, body map[string]interface{}, headerMap map[string]string) ([]byte, error) {
	statusCode, byteSlice, err := RequestByteSliceResult("DELETE", url, body, headerMap)
	text := string(byteSlice)
	if err != nil {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, err}
	} else if statusCode == 200 || statusCode == 202 || statusCode == 204 {
		return byteSlice, nil
	} else {
		return byteSlice, RequestError{url, statusCode, byteSlice, &text, nil}
	}
}
