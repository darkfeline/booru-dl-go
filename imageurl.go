package dl

import (
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// retrieveImageURL returns the image URL for the booru URL.
func retrieveImageURL(rawurl string) (string, error) {
	r, err := http.Get(rawurl)
	if err != nil {
		return "", errors.Wrapf(err, "getting %s", rawurl)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return "", errors.Errorf("download %s: HTTP %d %s", rawurl, r.StatusCode, r.Status)
	}
	i, err := findImageURL(rawurl, r.Body)
	if err != nil {
		return "", errors.Wrapf(err, "cannot find image URL for %s", rawurl)
	}
	return i, nil
}

// findImageURL returns the image URL for a booru page.
func findImageURL(rawurl string, r io.Reader) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", errors.Wrapf(err, "cannot parse document %s", rawurl)
	}
	switch u.Hostname() {
	case "chan.sankakucomplex.com":
		return findSankakuImageURL(u, d)
	case "danbooru.donmai.us":
		return findDanbooruImageURL(u, d)
	default:
		return "", errors.Errorf("unknown booru %s", rawurl)
	}
}

func findSankakuImageURL(u *url.URL, d *goquery.Document) (string, error) {
	src, ok := d.Find("#highres").Attr("href")
	if !ok {
		return "", errors.New("cannot find image URL")
	}
	rel, err := url.Parse(src)
	if err != nil {
		return "", errors.Wrap(err, "parsing img src")
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
	u, err := url.Parse(src)
	if err != nil {
		return "", errors.Wrap(err, "parsing img src")
	}
	return u.String(), nil
}
