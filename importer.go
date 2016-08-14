package grepurl

import (
	"bufio"
	"os"
	"sync"
)

type URLData struct {
	url  string
	tlas []string
}

type TLAData struct {
	id   uint32
	tlas []string
}

func RunImport(files []string) {
	// channels and global state
	raw_urls := make(chan string, 1000)
	urls_to_store := make(chan URLData, 1000)
	tlas_to_index := make(chan TLAData, 1000)
	var fileProc, splitProc sync.WaitGroup

	urlstore := NewMemoryURLStore()
	tlaindex := NewTLAIndex()

	for _, fn := range files {
		ff := fn
		fileProc.Add(1)
		go splitFile(ff, raw_urls, &fileProc)
	}
	splitProc.Add(3)
	go splitURLs(raw_urls, urls_to_store, &splitProc)
	go buildTLAs(urls_to_store, tlas_to_index, urlstore, &splitProc)
	go ingestTLAs(tlas_to_index, tlaindex, &splitProc)
	fileProc.Wait()
	close(raw_urls)
	splitProc.Wait()

	tlaindex.Print()
	/*for tla := range tlas_to_index {
		fmt.Println(tla.id, tla.tlas)
	}*/

	//all this needs to complete
}

func ingestTLAs(tlas chan TLAData, idx *TLAIndex, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range tlas {
		idx.Add(item)
	}
}

// Convert
func buildTLAs(urls_to_index chan URLData, tlas_to_index chan TLAData, urlstore URLStore, wg *sync.WaitGroup) {
	for item := range urls_to_index {
		tlas_to_index <- TLAData{
			urlstore.addURL(item.url, item.tlas), item.tlas,
		}
	}

	close(tlas_to_index)
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

func splitURLs(inp chan string, outp chan URLData, wg *sync.WaitGroup) {
	for line := range inp {
		outp <- URLData{line, SplitNgram(line)}
	}
	close(outp)
	wg.Done()
}

// shared state for proce
