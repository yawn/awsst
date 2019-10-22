package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {

	assert := assert.New(t)

	var test SectionSlice = []string{
		"2.2.2.2.2",
		"1",
		"2.2.2.2.2.1.2.4.6",
		"2.2.2.2",
		"1.1.1",
		"2.2.2.2.1",
		"1.2",
		"1.1",
	}

	test.Sort()

	assert.Equal(SectionSlice{
		"1",
		"1.1",
		"1.1.1",
		"1.2",
		"2.2.2.2",
		"2.2.2.2.1",
		"2.2.2.2.2",
		"2.2.2.2.2.1.2.4.6",
	}, test)

}
