package grepurl

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/danohuiginn/btree"
	"github.com/stretchr/testify/assert"
)

func cmp(a, b int) int {
	return a - b
}

func TestBTree(t *testing.T) {
	_ = btree.TreeNew(cmp)
}

func TestBuildSorted(t *testing.T) {
	// for now, this is just testing the builtin sort function
	// but it will need to work with memory-mapped sorted arrays
	var tar [30]int
	testslice := tar[:]
	for i, _ := range testslice {
		testslice[i] = rand.Intn(9999)
	}
	inp := build_sorted_array(testslice)
	for i := 0; i < len(inp)-1; i++ {
		if inp[i] > inp[i+1] {
			t.Fail()
		}
	}
	fmt.Println(inp)
}

func TestSimpleIntersect(t *testing.T) {
	// find the intersection of two
	one := []int{5, 2, 7, 4, 8, 98, 97, 39, 50, 40}
	two := []int{44, 2, 7, 3, 9}
	expected_inters := []int{2, 7}
	inters := intersect_two(one, two)
	fmt.Println(inters)
	if !reflect.DeepEqual(inters, expected_inters) {
		t.Fail()
	}
	_ = inters
}

func randomInts(length int) []int {
	output := make([]int, length)
	for i, _ := range output {
		output[i] = rand.Intn(99999)
	}
	return output[:]
}

func TestIntersect(t *testing.T) {
	// try intersections with long lists of random ints
	// it takes ~6.5 seconds to handle intersection of 2
	// slices of 10 million integers each
	lengths := []int{1, 100, 1000, 10000, 100000}
	for _, len_one := range lengths {
		fmt.Print("building random ints: one...")
		fmt.Println("done")
		one := randomInts(len_one)
		for _, len_two := range lengths {
			fmt.Print("building random ints: two")
			two := randomInts(len_two)
			fmt.Println("done")
			fmt.Println(len_one, len_two)
			_ = intersect_two(one, two)
		}
	}
}

func TestSplitNgram(t *testing.T) {
	// see what happens when we split an url
	url := "example.com"
	expected := []string{START_URL + "ex", "exa", "xam", "amp", "mpl", "ple", "le.", "e.c", ".co", "com", "om" + END_URL}
	result := SplitNgram(url)
	for i, exp := range expected {
		_, _, _ = i, exp, result
		//assert.Equal(t, result[i], exp, "ngram splitting non-match")
	}
}

func TestNgramsFromFile(t *testing.T) {
	ngrams_from_file("/home/src/golang/src/github.com/danohuiginn/grepurl/testdata.txt")
	res := urlmatches("tiv")
	assert.Equal(t, res, []string{"http://exple.tive.org/blarg/2015/09/20/bourne-aesthetic",
		"http://kamiel.creativechoice.org/2015/09/10/will-work-for-the-commons/"})
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
