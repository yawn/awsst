package report

import (
	"bytes"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Report struct {
	AllNew             map[string]*Element
	AllOld             map[string]*Element
	Changed            map[string]*Element
	ChecksumsIdentical bool
	Delta              uint
	Missing            map[string]*Element // delete or renamed
	New                map[string]*Element // added or renamed
	SortedSections     []string
	Timediff           time.Duration
}

func (r *Report) ToText() string {

	if r.ChecksumsIdentical {
		return "Identical"
	}

	buf := bytes.NewBuffer(nil)

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Section", "Type", "Old", "New"})

	// TODO: add left and right as elements to report, report on timestamps not diffs

	types := map[string]map[string]*Element{
		"ADDED/RENAMED":  r.New,
		"DELETE/RENAMED": r.Missing,
		"CHANGED":        r.Changed,
	}

	for _, n := range r.SortedSections {

		for t, collection := range types {

			_, ok := collection[n]

			if ok {

				var (
					oText string
					nText string
				)

				if oe := r.AllOld[n]; oe != nil {
					oText = oe.Text
				}

				if ne := r.AllNew[n]; ne != nil {
					nText = ne.Text
				}

				table.Append([]string{
					n,
					t,
					oText,
					nText,
				})

				break

			}

		}

	}

	table.Render()

	return buf.String()

}
