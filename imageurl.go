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
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// retrieveImageURL returns the image URL for the booru image URL.
func retrieveImageURL(rawurl string) (string, error) {
	r, err := http.Get(rawurl)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return "", fmt.Errorf("GET %s %s", rawurl, r.Status)
	}
	i, err := findImageURL(r.Request.URL, r.Body)
	if err != nil {
		return "", err
	}
	return i, nil
}

// findImageURL returns the image URL for a booru page.
func findImageURL(u *url.URL, r io.Reader) (string, error) {
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	switch u.Hostname() {
	case "chan.sankakucomplex.com":
		return findSankakuImageURL(u, d)
	case "danbooru.donmai.us":
		return findDanbooruImageURL(u, d)
	default:
		return "", fmt.Errorf("unknown booru %s", u)
	}
}

func findSankakuImageURL(u *url.URL, d *goquery.Document) (string, error) {
	src, ok := d.Find("#highres").Attr("href")
	if !ok {
		return "", errors.New("cannot find image URL")
	}
	rel, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	i := u.ResolveReference(rel)
	return i.String(), nil
}

func findDanbooruImageURL(u *url.URL, d *goquery.Document) (string, error) {
	src, ok := d.Find("#image-resize-link").Attr("href")
	if !ok {
		src, ok = d.Find("#image").Attr("src")
		if !ok {
			return "", errors.New("cannot find image URL")
		}
	}
	return src, nil
}
