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

// Package dl implements downloading images from boorus.
package dl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func requestImage(rawurl string) (*http.Response, error) {
	i, err := retrieveImageURL(rawurl)
	if err != nil {
		return nil, fmt.Errorf("get URL: %s", err)
	}
	r, err := httpGetWithRetry(i)
	if err != nil {
		return nil, fmt.Errorf("get image: %s", err)
	}
	return r, nil
}

// WriteImage downloads the image from the booru URL and writes it to
// the Writer.
func WriteImage(rawurl string, w io.Writer) error {
	r, err := requestImage(rawurl)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = io.Copy(w, r.Body)
	return err
}

// Download downloads the image from the booru URL to the path.
func Download(rawurl string, p string) error {
	r, err := requestImage(rawurl)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return writeToFile(r.Body, p)
}

// DownloadToDir downloads the image from the booru URL to the
// directory using a default filename.
func DownloadToDir(rawurl string, dir string) error {
	r, err := requestImage(rawurl)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	f := urlFilename(r.Request.URL)
	fp := filepath.Join(dir, f)
	return writeToFile(r.Body, fp)
}

// urlFilename returns the filename to be used for saving the URL.
func urlFilename(u *url.URL) string {
	return path.Base(u.Path)
}

func writeToFile(r io.Reader, p string) error {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err2 := f.Close(); err == nil {
		err = err2
	}
	return err
}
