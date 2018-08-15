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
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// RetrieveImage retrieves the image from the given URL.  This
// function does magic to handle links to image boorus.  You must
// close ImageData.Data if there is no error.
func RetrieveImage(rawurl string) (*ImageData, error) {
	i, err := retrieveImageURL(rawurl)
	if err != nil {
		return nil, fmt.Errorf("get image URL for %s: %s", rawurl, err)
	}
	r, err := httpGetWithRetry(i)
	if err != nil {
		return nil, fmt.Errorf("get image for %s: %s", rawurl, err)
	}
	d := &ImageData{
		Data:    r.Body,
		FileURL: r.Request.URL,
	}
	return d, nil
}

// ImageData contains the retrieved image data and metadata.
type ImageData struct {
	// Data is the image data.  You must close this.
	Data io.ReadCloser
	// URL is the URL of the image file.
	FileURL *url.URL
}

// WriteImage downloads the image from the booru URL and writes it to
// the Writer.
func WriteImage(rawurl string, w io.Writer) error {
	d, err := RetrieveImage(rawurl)
	if err != nil {
		return err
	}
	defer d.Data.Close()
	_, err = io.Copy(w, d.Data)
	return err
}

// Download downloads the image from the booru URL to the path.
func Download(rawurl string, p string) error {
	d, err := RetrieveImage(rawurl)
	if err != nil {
		return err
	}
	defer d.Data.Close()
	return writeToFile(d.Data, p)
}

// DownloadToDir downloads the image from the booru URL to the
// directory using a default filename.
func DownloadToDir(rawurl string, dir string) error {
	d, err := RetrieveImage(rawurl)
	if err != nil {
		return err
	}
	defer d.Data.Close()
	f := urlFilename(d.FileURL)
	fp := filepath.Join(dir, f)
	return writeToFile(d.Data, fp)
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
