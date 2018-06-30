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
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestFindImageURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		url, f, exp string
	}{
		{"https://chan.sankakucomplex.com/post/show/6932916", "sankaku-full",
			"https://cs.sankakucomplex.com/data/ce/5b/ce5b80eb7bfe73f1ec4a8ff348e286aa.jpg?e=1528064163&m=N1SxFIx2X1MQWmSBHyakqw"},
		{"https://chan.sankakucomplex.com/post/show/6935519", "sankaku-sample",
			"https://cs.sankakucomplex.com/data/93/7f/937f92fd467786d17bebcef825b4c547.jpg?e=1528066908&m=YRmTT_4g4Rd082EhiZ4Hfw"},
		{"https://danbooru.donmai.us/posts/2774553", "danbooru-full",
			"https://hijiribe.donmai.us/data/__kunikida_hanamaru_love_live_and_love_live_sunshine_drawn_by_mignon__1fc1c222db895e7e0ee5da78886ee5f0.jpg"},
		{"https://danbooru.donmai.us/posts/3146000", "danbooru-sample",
			"https://danbooru.donmai.us/data/__enoch_soul_worker_drawn_by_noria__7084f5bc4fb04da9974dc50546737c62.jpg"},
	}
	for _, c := range cases {
		f, err := os.Open(filepath.Join("testdata", fmt.Sprintf("%s.html", c.f)))
		if err != nil {
			t.Errorf("Error opening test file: %s", err)
			continue
		}
		u, err := url.Parse(c.url)
		if err != nil {
			t.Fatalf("Error parsing URL %s: %s", c.url, err)
		}
		got, err := findImageURL(u, f)
		if err != nil {
			t.Errorf("findImageURL return error: %s", err)
			continue
		}
		if got != c.exp {
			t.Errorf("For %s, expected %s, got %s", c.f, c.exp, got)
		}
	}
}
