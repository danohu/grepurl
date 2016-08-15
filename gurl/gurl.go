package main

import (
	"fmt"

	grepurl "github.com/danohuiginn/grepurl"
)

func main() {
	files := []string{"/tmp/urlsample.txt"}
	urlstore, trigrams := grepurl.RunImport(files)
	fmt.Println("import complete")
	query := ".*.gov"
	ch := make(chan string)
	go grepurl.RunQuery(query, trigrams, urlstore, ch)
	for result := range ch {
		fmt.Println(result)
	}
}
