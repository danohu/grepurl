package grepurl

import(
	"fmt"
	"sort"
)

func build_sorted_array(unsorted []int) []int {
	// this should also be removing duplicates?
	sort.Ints(unsorted)
	return unsorted
}

func intersect_two_slow(one []int, two []int) []int{
	inters := []int {}
	one = build_sorted_array(one)
	two = build_sorted_array(two)
	operations := 0
	for i_one, _ := range one{
		for i_two, _ := range two {
			operations += 1
			if(one[i_one] == two[i_two]){
				inters = append(inters, one[i_one])
			}
		}
	}
	fmt.Println(operations, "operations")
	return inters
}


func intersect_two(one []int, two []int) []int{
	inters := []int {}
	one = build_sorted_array(one)
	two = build_sorted_array(two)
	operations := 0
	i_two := 0
	for i_one, _ := range one{
		for i_two < len(two){
			operations += 1
			if(two[i_two] == one[i_one]){
				inters = append(inters, one[i_one])
				break
			}
			if two[i_two] > one[i_one]{
				break
			}
			i_two += 1
		}
	}
	fmt.Println(operations, "operations")
	return inters
}


func main(){
	fmt.Println("Hello world")
}
