package grepurl

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/* */

var MonthIndexes = map[string]string{
	"2016-06": "2016-26",
	"2016-05": "2016-22",
}

const bucket_name = "commoncrawl"

//XXX month should look like "2016-06"
// beware tat
// -- just getting the latest month for onw
func ListIndexFiles(month string) []string {

	indexfiles := []string{}

	monthindex := MonthIndexes[month]
	file_name := "/crawl-data/CC-MAIN-" + monthindex + "/wat.paths.gz"
	fmt.Println(file_name)

	scanner := s3Gzip(bucket_name, file_name)
	for scanner.Scan() {
		indexfiles = append(indexfiles, scanner.Text())
	}
	return indexfiles
}

/* This is assuming basic common_crawl stuff

 */
func URLsFromWAT(path string, ch chan string) {
	scanner := s3Gzip(bucket_name, path)
	/*f, _ := os.Open("testdata/wat.gz")
	gr, _ := gzip.NewReader(f)
	scanner := bufio.NewScanner(gr)
	*/
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, "WARC-Target-URI: ") {
			url := strings.SplitN(line, ": ", 2)
			ch <- url[1]
		}
	}
	close(ch)
}

/*
Return the number of URLs uploaded,
or -1 if we didn't upload anything
(e.g. because we avoided replacing an existing file)
*/
func UploadUrls(path string, clobber bool) int {
	return -1
}

// Take a gzipped s3 file, and return a scanner over it
// Internally the file is downloaded to a temporary location,
// unzipped and then removed
func s3Gzip(bucket_name string, file_name string) *bufio.Scanner {
	file, _ := ioutil.TempFile("", "")
	defer file.Close()
	defer os.Remove(file.Name())

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(file_name),
		})

	gr, _ := gzip.NewReader(file)
	return bufio.NewScanner(gr)
}
