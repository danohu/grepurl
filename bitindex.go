package grepurl

import (
	"fmt"

	roaring "github.com/RoaringBitmap/roaring"
)

type Trigram struct {
	letters string
	bitmap  roaring.Bitmap
}

func (ng *Trigram) GetSet(ch chan uint32) {
	for _, el := range ng.bitmap.ToArray() {
		ch <- el
	}
	close(ch)
}

type TrigramIndex struct {
	cardinality int
	letters     map[string]*roaring.Bitmap
}

func NewTrigramIndex() *TrigramIndex {
	return &TrigramIndex{
		cardinality: 8388608, // 1 Megabyte in bits
		letters:     make(map[string]*roaring.Bitmap)}
}

func (ind *TrigramIndex) Add(data TrigramData) {
	for _, trigram := range data.trigrams {
		_, exists := ind.letters[trigram]
		if !exists {
			ind.letters[trigram] = roaring.NewBitmap()
		}
		ind.letters[trigram].Add(data.id)
	}
}

func (ind *TrigramIndex) Print() {
	for k, v := range ind.letters {
		fmt.Println(k, v.String())
	}
}

func MiniBit() {
	ch := make(chan uint32)
	roar := roaring.BitmapOf(14, 1000, 99902)
	ng := &Trigram{letters: "wtf", bitmap: *roar}
	go ng.GetSet(ch)
	for el := range ch {
		fmt.Println(el)
	}

}
