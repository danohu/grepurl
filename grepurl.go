package grepurl

import(
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	)

//Map id -> URL
var urlcodes map[int]string
var nexturlcode int

// map TLA --> list of items
var occurances = make(map[string][]int)


func _init(){
	//set up the global vars
	urlcodes = make(map[int]string)
	nexturlcode = 0
}

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

func SplitNgram(url string) []string{
	// Split an url into three-character chunks
	results := []string{}
	for i := 0; i < len(url) - 2; i++ {
		results = append(results, url[i:i+3])
	}
	return results
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

func ngrams_from_file(filepath string){
	_init()

	file, err := os.Open(filepath)
	if err != nil{
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := scanner.Text()
		urlcodes[nexturlcode] = line
		for _, ng := range SplitNgram(line){
			occurances[ng] = append(occurances[ng], nexturlcode)
			fmt.Println(ng)
			}
		nexturlcode++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func urlmatches(tla string) []string{
	res := []string{}
	for _, i := range occurances[tla]{
		res = append(res, urlcodes[i])
	}
	return res
}


func main(){
	fmt.Println("Hello world")
}


