package grepurl

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
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
		actual, _ := ms.getURL(idx)
		assert.Equal(t, el, actual)
		_ = idx
	}
	actual, _ := ms.getURL(0x01)
	assert.Equal(t, "bbb", actual)

	_, err := ms.getURL(0x99)
	assert.NotEqual(t, nil, err)

	_, err = ms.getURL(0x04)
	assert.NotEqual(t, nil, err)

}

func TestRunImport(t *testing.T) {
	files := []string{"testdata/urls.txt"}
	RunImport(files)
}

/*func TestURLsFromWat(t *testing.T) {
	ch := make(chan string)
	go URLsFromWAT("", ch)
	for url := range ch {
		fmt.Println(url)
	}
}*/

func TestUploadURLs(t *testing.T) {
	res := UploadURLs("", false)
	assert.Equal(t, res, -1)
}

func TestGenerateS3Path(t *testing.T) {
	source := "crawl-data/CC-MAIN-2016-26/segments/1466783408840.13/wat/CC-MAIN-20160624155008-00193-ip-10-164-35-72.ec2.internal.warc.wat.gz"
	expected := "urls/2016-26/CC-MAIN-20160624155008-00193-ip-10-164-35-72.urls.gz"
	actual := GenerateS3Path(source)
	assert.Equal(t, expected, actual)
}

func TestS3PathExists(t *testing.T) {
	exists := "crawl-data/CC-MAIN-2016-26/segments/1466783408840.13/wat/CC-MAIN-20160624155008-00193-ip-10-164-35-72.ec2.internal.warc.wat.gz"
	does_not_exist := "crawl-data/CC-MAIN-2016-26/segments/1466783408840.13/wat/CC-MAIN-20160624155008-00193-ip-10-164-35-72.ec2.internal.BLAH_BLAH_BLAH.warc.wat.gz"
	assert.Equal(t, true,
		S3PathExists("commoncrawl", exists, "us-east-1"))
	assert.Equal(t, false,
		S3PathExists("commoncrawl", does_not_exist, "us-east-1"))
}

func setUp() (URLStore, *TrigramIndex) {
	files := []string{"testdata/urls.txt"}
	return RunImport(files)
}

func TestRetrieveUrl(t *testing.T) {
	urlstore, trigrams := setUp()
	query, expected := "fabians.*ELect...l", "http://www.fabians.org.uk/under-corbyns-ELectoral-plan-prospects-for-victory-look-bleak/"
	ch := make(chan string)
	go RunQuery(query, trigrams, urlstore, ch)
	actual := <-ch
	assert.Equal(t, expected, actual)
}

func TestRegex(t *testing.T) {
	var sample_queries = map[string][]uint32{
		"tiv":                []uint32{2, 3},
		"(tiv|non-matching)": []uint32{2, 3},
	}
	_, tgindex := setUp()

	for rgx, matches_expected := range sample_queries {
		qry := BuildQueryObject(rgx)
		ch := make(chan uint32)
		go ApplyQuery(qry, tgindex, ch)
		_ = matches_expected
		matches_actual := make([]uint32, 0)
		for i := range ch {
			matches_actual = append(matches_actual, i)
		}
		assert.Equal(t, matches_expected, matches_actual)

	}
}

// with empty bitmaps, we should ever get any matches
func TestEmptyRegex(t *testing.T) {
	regexes := []string{"tiv", "[abc]def", "this_is_just_a_long_word2"}
	_, tgindex := RunImport([]string{})

	for _, rgx := range regexes {
		qry := BuildQueryObject(rgx)
		roar := RoaringQuery(qry, tgindex)
		assert.Equal(t, roaring.NewBitmap(), roar)
	}
}
