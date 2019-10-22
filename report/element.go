package report

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/yawn/awsst/sort"
)

type Element struct {
	Checksum []byte
	Number   string
	Text     string
}

type TOS struct {
	ArchiveURL string
	Checksum   []byte
	Elements   []*Element
	Timestamp  time.Time
}

func New(number, text string) *Element {

	h := sha256.New()

	h.Write([]byte(number))
	h.Write([]byte(text))

	return &Element{
		Checksum: h.Sum(nil),
		Number:   number,
		Text:     text,
	}

}

func (t *TOS) Compare(new *TOS) *Report {

	report := &Report{
		AllNew:             make(map[string]*Element),
		AllOld:             make(map[string]*Element),
		Changed:            make(map[string]*Element),
		ChecksumsIdentical: bytes.Equal(t.Checksum, new.Checksum),
		Missing:            make(map[string]*Element),
		New:                make(map[string]*Element),
		Timediff:           t.Timestamp.Sub(new.Timestamp),
	}

	var (
		all    = make(map[string]struct{})
		sorted sort.SectionSlice
	)

	for _, en := range new.Elements {
		all[en.Number] = struct{}{}
		report.AllNew[en.Number] = en
	}

	for _, eo := range t.Elements {
		all[eo.Number] = struct{}{}
		report.AllOld[eo.Number] = eo
	}

	for n := range all {
		sorted = append(sorted, n)
	}

	sorted.Sort()

	report.SortedSections = sorted

	// check for new
	for _, en := range new.Elements {

		_, ok := report.AllOld[en.Number]

		if !ok {
			report.New[en.Number] = en
			report.Delta++
		}

	}

	// check for missing or changed
	for _, eo := range t.Elements {

		en, ok := report.AllNew[eo.Number]

		if !ok {
			report.Missing[eo.Number] = eo
			report.Delta++
		} else if !bytes.Equal(eo.Checksum, en.Checksum) {
			report.Changed[eo.Number] = eo
			report.Delta++
		}

	}

	return report

}
