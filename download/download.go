package download

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	_url "net/url"
	"time"

	"github.com/pkg/errors"
)

type Callback func(*Period, []byte) error

type Downloader struct {
	url            string
	SkipDuplicates bool
}

type Period struct {
	Day   uint
	Month uint
	Year  uint
}

func (p *Period) String() string {

	if p == nil {
		return "current"
	} else {
		return fmt.Sprintf("%4d-%02d-%02d", p.Year, p.Month, p.Day)
	}

}

func New(url string) (*Downloader, error) {

	u, err := _url.Parse(url)

	if err != nil {
		return nil, err
	}

	return &Downloader{
		url: u.Hostname() + u.RequestURI(),
	}, nil

}

func (d *Downloader) Download(cb Callback, periods []*Period) error {

	sums := make(map[string]*Period)

	for _, p := range periods {

		var url string

		if p == nil {
			url = fmt.Sprintf("https://%s", d.url)
		} else {

			url = fmt.Sprintf("https://web.archive.org/web/%4d%02d%02d000000/%s",
				p.Year,
				p.Month,
				p.Day,
				d.url)

		}

		log.Printf("downloading %q", url)

		res, err := http.Get(url)

		if err != nil {
			return errors.Wrapf(err, "failed to download %q", url)
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("status code error: %d %q", res.StatusCode, res.Status)
		}

		var (
			buf = bytes.NewBuffer(nil)
			h   = sha256.New()
		)

		if _, err := io.Copy(buf, io.TeeReader(res.Body, h)); err != nil {
			return errors.Wrapf(err, "failed to read body of %q", url)
		}

		if d.SkipDuplicates {

			sum := hex.EncodeToString(h.Sum(nil))

			if period, ok := sums[sum]; ok {

				log.Printf("skipping %s, it's a duplicate of %s (%s)",
					url,
					period,
					sum,
				)

				continue

			}

			sums[sum] = p

		}

		if cb != nil {

			if p == nil {

				now := time.Now()

				p = &Period{
					Day:   uint(now.Day()),
					Month: uint(now.Month()),
					Year:  uint(now.Year()),
				}

			}

			if err := cb(p, buf.Bytes()); err != nil {
				return err
			}

		}

		res.Body.Close()

	}

	return nil

}
