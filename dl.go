/*
Package dl implements downloading images from boorus.
*/
package dl

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

// Download downloads the image from the booru URL to the path.
func Download(rawurl string, p string) error {
	i, err := retrieveImageURL(rawurl)
	if err != nil {
		return errors.Wrapf(err, "cannot find image URL for %s", rawurl)
	}
	err = saveHTTP(i, p)
	if err != nil {
		return errors.Wrapf(err, "save %s to %s", i, p)
	}
	return nil
}

// DownloadToDir downloads the image from the booru URL to the
// directory using a default filename.
func DownloadToDir(rawurl string, dir string) error {
	i, err := retrieveImageURL(rawurl)
	if err != nil {
		return errors.Wrapf(err, "cannot find image URL for %s", rawurl)
	}
	f, err := urlFilename(i)
	if err != nil {
		return errors.Wrapf(err, "cannot find filename for %s", i)
	}
	fp := filepath.Join(dir, f)
	err = saveHTTP(i, fp)
	if err != nil {
		return errors.Wrapf(err, "save %s to %s", i, fp)
	}
	return nil
}

// urlFilename returns the filename to be used for saving the URL.
func urlFilename(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", errors.Wrap(err, "parse URL filename")
	}
	return path.Base(u.Path), nil
}

func saveHTTP(rawurl string, p string) error {
	r, err := http.Get(rawurl)
	if err != nil {
		return errors.Wrapf(err, "getting %s", rawurl)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return errors.Errorf("download %s: HTTP %d %s", rawurl, r.StatusCode, r.Status)
	}
	return writeToFile(r.Body, p)
}

func writeToFile(r io.Reader, p string) error {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
