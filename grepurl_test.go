package grepurl

import (
	"testing"
	"fmt"
	"math/rand"
	"reflect"
)

func TestBuildSorted(t *testing.T){
	// for now, this is just testing the builtin sort function
	// but it will need to work with memory-mapped sorted arrays	
	var tar [30]int
	testslice := tar[:]
	for i,_ := range testslice {
		testslice[i] = rand.Intn(9999)
	}
	inp := build_sorted_array(testslice)
	for i := 0; i<len(inp)-1;i++ {
		if(inp[i] > inp[i+1]){
			t.Fail()	
		}
	}
	fmt.Println(inp)
}

func TestSimpleIntersect(t *testing.T){
	// find the intersection of two
	one := []int{5,2,7,4,8,98,97,39,50,40}
	two := []int{44,2,7,3,9}
	expected_inters := []int{2,7}
	inters := intersect_two(one, two)
	fmt.Println(inters)
	if(!reflect.DeepEqual(inters, expected_inters)){
		t.Fail()
	}
	_ = inters
}

func randomInts(length int) []int{
	output := make([]int, length)
	for i,_ := range output{
		output[i] = rand.Intn(99999)
	}
	return output[:]
}

func TestIntersect(t *testing.T){
	// try intersections with long lists of random ints
	lengths := [] int {1,100,1000,10000,100000} //,100000000}
	for _, len_one := range lengths {
		one := randomInts(len_one)
		for _, len_two := range lengths {
			two := randomInts(len_two)
			fmt.Println(len_one, len_two)
			_ = intersect_two(one,two)
			/*res_slow := res_fast
			//res_slow := intersect_two(one,two)
			if(!reflect.DeepEqual(res_fast, res_slow)){
				fmt.Println(res_fast)
				fmt.Println(res_slow)
				t.Fail()
			}*/

		}
	}
}
