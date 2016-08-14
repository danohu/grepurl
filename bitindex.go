package grepurl

import (
	"fmt"

	roaring "github.com/RoaringBitmap/roaring"
)

type NGram struct {
	letters string
	bitmap  roaring.Bitmap
}

func (ng *NGram) GetSet(ch chan uint32) {
	for _, el := range ng.bitmap.ToArray() {
		ch <- el
	}
	close(ch)
}

type TLAIndex struct {
	cardinality int
	letters     map[string]*roaring.Bitmap
}

func NewTLAIndex() *TLAIndex {
	return &TLAIndex{
		cardinality: 8388608, // 1 Megabyte in bits
		letters:     make(map[string]*roaring.Bitmap)}
}

func (ind *TLAIndex) Add(data TLAData) {
	for _, tla := range data.tlas {
		_, exists := ind.letters[tla]
		if !exists {
			ind.letters[tla] = roaring.NewBitmap()
		}
		ind.letters[tla].Add(data.id)
	}
}

func (ind *TLAIndex) Print() {
	for k, v := range ind.letters {
		fmt.Println(k, v.String())
	}
}

func MiniBit() {
	ch := make(chan uint32)
	roar := roaring.BitmapOf(14, 1000, 99902)
	ng := &NGram{letters: "wtf", bitmap: *roar}
	go ng.GetSet(ch)
	for el := range ch {
		fmt.Println(el)
	}

}
