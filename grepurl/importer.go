package grepurl

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const START_URL = "\x02"
const END_URL = "\x03"

type URLData struct {
	url      string
	trigrams []string
}

type TrigramData struct {
	id       uint32
	trigrams []string
}

func RunImport(files []string) (URLStore, *TrigramIndex) {
	// channels and global state
	raw_urls := make(chan string, 1000)
	urls_to_store := make(chan URLData, 1000)
	trigrams_to_index := make(chan TrigramData, 1000)
	var fileProc, splitProc sync.WaitGroup

	urlstore := NewMemoryURLStore()
	trigramindex := NewTrigramIndex()

	for _, fn := range files {
		ff := fn
		fileProc.Add(1)
		go splitFile(ff, raw_urls, &fileProc)
	}
	splitProc.Add(3)
	go splitURLs(raw_urls, urls_to_store, &splitProc)
	go buildTrigrams(urls_to_store, trigrams_to_index, urlstore, &splitProc)
	go ingestTrigrams(trigrams_to_index, trigramindex, &splitProc)
	fileProc.Wait()
	close(raw_urls)
	splitProc.Wait()

	return urlstore, trigramindex
}

func ingestTrigrams(trigrams chan TrigramData, idx *TrigramIndex, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range trigrams {
		idx.Add(item)
	}
}

// Convert
func buildTrigrams(urls_to_index chan URLData, trigrams_to_index chan TrigramData, urlstore URLStore, wg *sync.WaitGroup) {
	for item := range urls_to_index {
		trigrams_to_index <- TrigramData{
			urlstore.addURL(item.url, item.trigrams), item.trigrams,
		}
	}

	close(trigrams_to_index)
	wg.Done()
}

// push lines in a file out into a channel
func splitFile(fn string, raw_urls chan string, wg *sync.WaitGroup) {
	file, _ := os.Open(fn)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw_urls <- scanner.Text()
	}
	wg.Done()
}

func SplitTrigram(url string) []string {
	// Split an url into three-character chunks
	results := []string{}
	url = PrepareUrl(url)
	url = START_URL + url + END_URL // XX should this go in PrepareUrl?
	for i := 0; i < len(url)-2; i++ {
		results = append(results, url[i:i+3])
	}
	return results
}

//Prepare an URL for the trigram index
//This consists of:
// - lowercasing everything
// - optionally adding start and end anchors
// - applying url-encoding where applicable
// - stripping out any remaining non-allowed characters
func PrepareUrl(raw string) string {
	prepared := strings.ToLower(raw)
	return prepared
}

func splitURLs(inp chan string, outp chan URLData, wg *sync.WaitGroup) {
	for line := range inp {
		outp <- URLData{line, SplitTrigram(line)}
	}
	close(outp)
	wg.Done()
}

// shared state for proce
