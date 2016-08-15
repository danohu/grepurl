package grepurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTrigram(t *testing.T) {
	// see what happens when we split an url
	url := "example.com"
	expected := []string{START_URL + "ex", "exa", "xam", "amp", "mpl", "ple", "le.", "e.c", ".co", "com", "om" + END_URL}
	result := SplitTrigram(url)
	for i, exp := range expected {
		_, _, _ = i, exp, result
		assert.Equal(t, result[i], exp, "ngram splitting non-match")
	}
}

func TestURLStore(t *testing.T) {
	ms := NewMemoryURLStore()
	for i, el := range []string{"aa", "bbb", "c", "d"} {
		idx := ms.addURL(el, []string{})
		assert.Equal(t, uint32(i)+0x01, ms.getCardinality())
		assert.Equal(t, el, ms.getURL(idx))
		_ = idx
	}
	assert.Equal(t, "bbb", ms.getURL(0x01))
}

func TestRunImport(t *testing.T) {
	files := []string{"testdata.txt"}
	RunImport(files)
}

func TestRegex(t *testing.T) {
	BuildQuery("([aA]bc|defgh)")
}
