// Copyright (C) 2018 Allen Li
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

package dl

import (
	"fmt"
	"net/http"
	"time"
)

const httpRetries = 5
const httpRetryInterval = 5 * time.Second

func httpGetWithRetry(rawurl string) (*http.Response, error) {
	var err error
	for i := 0; i < httpRetries; i++ {
		var r *http.Response
		r, err = httpGet(rawurl)
		if err != nil {
			if isTemporary(err) {
				time.Sleep(httpRetryInterval)
				continue
			}
			return nil, err
		}
		return r, nil
	}
	return nil, err
}

var _ error = tooManyRequestsErr{}
var _ temporaryErr = tooManyRequestsErr{}

type tooManyRequestsErr struct{}

func (tooManyRequestsErr) Error() string {
	return "too many requests"
}

func (tooManyRequestsErr) Temporary() bool {
	return true
}

type temporaryErr interface {
	Temporary() bool
}

func isTemporary(e error) bool {
	te, ok := e.(temporaryErr)
	if !ok {
		return false
	}
	return te.Temporary()
}

// httpGet performs an HTTP GET request.
//
// When err is nil, resp always contains a non-nil resp.Body. Caller
// should close resp.Body when done reading from it.
//
// err may be a temporaryErr.
func httpGet(rawurl string) (resp *http.Response, err error) {
	r, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer func(r *http.Response) {
		if resp == nil {
			r.Body.Close()
		}
	}(r)
	switch r.StatusCode {
	case http.StatusTooManyRequests:
		return nil, tooManyRequestsErr{}
	case http.StatusOK:
		return r, nil
	default:
		return nil, fmt.Errorf("GET %s %s", rawurl, r.Status)
	}
}
