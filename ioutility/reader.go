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

package ioutility

import (
	"io"
)

// Read from stream, keep reading until no data. However, EOF is not read and is a signal to close the input stream
func ReadText(reader io.Reader, bufferSize int) (string, int, error) {
	buffer := make([]byte, bufferSize)

	totalLength := 0
	byteSlice := make([]byte, 0)
	for {
		// read a chunk
		n, err := reader.Read(buffer)
		if err == io.EOF {
			byteSlice = append(byteSlice, buffer[0:n]...)
			return string(byteSlice), totalLength, err
		} else if err != nil {
			return string(byteSlice), totalLength, err
		} else if n == 0 {
			return string(byteSlice), totalLength, nil
		} else {
			byteSlice = append(byteSlice, buffer[0:n]...)
			totalLength += n
			if n < bufferSize {
				return string(byteSlice), totalLength, nil
			}
		}
	}
}
