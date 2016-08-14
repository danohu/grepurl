package grepurl

import (
	"fmt"

	bitset "github.com/willf/bitset"
)

type ONGram struct {
	letters string
	bitmap  bitset.BitSet
}

func (ng *ONGram) GetSet(ch chan uint) {
	var i uint
	for e := true; e == true; {
		i, e = ng.bitmap.NextSet(i)
		if e {
			ch <- i
		}
		i += 1
	}
	close(ch)

}

func GetSetBits(bs *bitset.BitSet, ch chan uint) {
	var i uint
	for e := true; e == true; {
		i, e = bs.NextSet(i)
		if e {
			ch <- i
		}
		i += 1
	}
	close(ch)
}

func MiniBit() {
	var a, b bitset.BitSet
	ch := make(chan uint)
	a.Set(0).Set(10).Set(15)
	b.Set(10).Set(11).Set(0).Set(15)
	c := a.Intersection(&b)
	ng := &ONGram{letters: "omg", bitmap: *c}

	go ng.GetSet(ch)
	for i := range ch {
		fmt.Println(i)
	}

}
