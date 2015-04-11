package grepurl

import(
	"fmt"
	"sort"
)


func build_sorted_array(unsorted []int) []int {
	// this should also be removing duplicates?
	filtered := [] int {}
	sort.Ints(unsorted)

	filtered = append(filtered, unsorted[0])
	
	for i := 1; i < len(unsorted) ; i++  {
		if(unsorted[i] != unsorted[i-1]){
			filtered = append(filtered, unsorted[i])
		}
	}
	return filtered
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
	for _, el_one := range one{
		for i_two < len(two){
			operations += 1
			if(two[i_two] == el_one){
				inters = append(inters, el_one)
				break
			}
			if two[i_two] > el_one{
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
