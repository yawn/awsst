package aws

import (
	"crypto/sha256"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/yawn/awsst/report"
)

const URL = "aws.amazon.com/service-terms"

var (
	rSection   = regexp.MustCompile(`(?s)^([\d\.]+)(?:\s+)(.+)`)
	rTimestamp = regexp.MustCompile(`(?i)last updated: (.+)`)
)

type Terms struct {
	report.TOS
	counts []int64
	depth  int
}

func New() *Terms {
	return new(Terms)
}

func (t *Terms) Parse(r io.Reader) error {

	var (
		h       = sha256.New()
		loopErr error
	)

	doc, err := goquery.NewDocumentFromReader(io.TeeReader(r, h))

	if err != nil {
		return errors.Wrapf(err, "failed to open document")
	}

	// parse timestamp
	doc.Find("h1, h2, h3, h4, h5").EachWithBreak(func(i int, s *goquery.Selection) bool {

		text := s.Text()

		if rTimestamp.MatchString(text) {

			date, err := time.Parse("January 2, 2006", rTimestamp.FindStringSubmatch(text)[1])

			if err != nil {
				loopErr = errors.Wrapf(err, "failed to parse date from %q", text)
				return false
			}

			t.ArchiveURL = fmt.Sprintf("https://web.archive.org/web/%4d%02d01000000/aws.amazon.com/service-terms", date.Year(), date.Month())
			t.Timestamp = date

		}

		return true

	})

	if loopErr != nil {
		return loopErr
	}

	// parse sections and subsections
	doc.Find("h3, h4, p").EachWithBreak(func(i int, s *goquery.Selection) bool {

		text := s.Text()
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\u00a0", " ") // nbsp

		if rSection.MatchString(text) {

			if err := t.parse(text); err != nil {
				loopErr = errors.Wrapf(err, "failed to parse text %q", text)
				return false
			}

		}

		return true

	})

	if loopErr != nil {
		return loopErr
	}

	t.Checksum = h.Sum(nil)

	return nil

}

func (t *Terms) parse(text string) error {

	var (
		parts    = rSection.FindStringSubmatch(text)
		section  = strings.TrimSuffix(parts[1], ".")
		elements = strings.Split(section, ".")
	)

	// how deep are we in?
	t.depth = len(elements) - 1

	// if we ascend, drop the lower level (e.g. 2.3.2, 2.3.3, 2.4 -> drop level 3)
	if t.depth < len(t.counts) {
		t.counts = t.counts[:t.depth+1]
	}

	// if it's a new depth, add one level
	if len(elements) > len(t.counts) {
		t.counts = append(t.counts, 0)
	}

	// check for unexpected gaps, that's all we can do really
	for idx, e := range elements {

		n, err := strconv.ParseInt(e, 10, 16)

		if err != nil {
			return err
		}

		if n > t.counts[idx]+1 {
			return fmt.Errorf("unexpected section %s", parts[1])
		}

		t.counts[idx] = n

	}

	// determine next section and check if it occurs in our text (due to malformed HTML)
	// TODO: also determine sections the levels above
	next := make([]string, len(t.counts))

	for idx, c := range t.counts {

		var add int64

		if idx == len(t.counts)-1 {
			add = 1
		}

		next[idx] = strconv.Itoa(int(c + add))

	}

	idx := strings.Index(text, strings.Join(next, "."))

	if idx > 0 {
		t.Elements = append(t.Elements, report.New(section, text[0:idx]))
		return t.parse(text[idx:])
	}

	t.Elements = append(t.Elements, report.New(section, text))

	return nil

}
