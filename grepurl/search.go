package grepurl

import (
	//index "github.com/google/codesearch/index"
	"fmt"
	"regexp"
	syntax "regexp/syntax"

	"github.com/RoaringBitmap/roaring"
	index "github.com/google/codesearch/index"
)

func printQuery(info *index.Query, depth string) {
	fmt.Println(depth, info.Trigram, info.Sub, info.Op)
	for _, s := range info.Sub {
		printQuery(s, depth+"*")
	}
}

func BuildQueryObject(inp string) *index.Query {
	rxp, _ := syntax.Parse(inp, syntax.Perl)
	info := index.RegexpQuery(rxp)
	printQuery(info, "")
	return info
}

func RunQuery(inp string, tgindex *TrigramIndex, urlstore URLStore, out chan string) {
	ids_in := make(chan uint32)

	// this is the full regexp
	// So it gets the unmodified string, so it can look at case etc
	// XXX needs handling of illegal URLs
	pattern := regexp.MustCompile(inp)

	trigram_ready_url := PrepareUrl(inp)
	qry := BuildQueryObject(trigram_ready_url)
	go ApplyQuery(qry, tgindex, ids_in)
	for id := range ids_in {
		url, err := urlstore.getURL(id)
		if err == nil {
			if pattern.MatchString(url) {
				out <- url
			}
		}
	}
	close(out)
}

// Given a codesearch query object and a collection of bitmaps
// showing which items match each trigram, push into a channel
// the ids of items which (probably) match the trigram requirements
func ApplyQuery(qry *index.Query, tgindex *TrigramIndex, out chan uint32) {
	roar := RoaringQuery(qry, tgindex)
	for _, j := range roar.ToArray() {
		out <- j
	}
	close(out)
}

func RoaringQuery(qry *index.Query, tgindex *TrigramIndex) *roaring.Bitmap {

	var bitmaps []*roaring.Bitmap

	// get all direct trigrams
	for _, trigram := range qry.Trigram {
		bmp, exists := tgindex.letters[trigram]
		if !exists {
			bmp = roaring.NewBitmap() // empty bitmap
		}
		bitmaps = append(bitmaps, bmp)
	}

	// get sub-queries
	for _, sub := range qry.Sub {
		bitmaps = append(bitmaps, RoaringQuery(sub, tgindex))
	}

	if index.QAnd == qry.Op {
		return roaring.FastAnd(bitmaps...)
	} else if index.QOr == qry.Op {
		return roaring.FastOr(bitmaps...)
	} else if index.QAll == qry.Op {
		// XXX inefficiently filling a bitmap with ones
		// nb this risks giving us indices that don't exist
		roar := roaring.BitmapOf(uint32(tgindex.cardinality))
		roar.AddRange(0, uint64(tgindex.cardinality))
		return roar
		//return roaring.NewBitmap()
	} else if index.QNone == qry.Op {
		return roaring.NewBitmap()
	}
	panic("incomprehensible query object")
}
