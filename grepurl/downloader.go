package grepurl

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/* */

//XXX month is curretly ignored
// -- just getting the latest month for onw
func ListIndexFiles(month string) []string {

	indexfiles := []string{}

	file, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	bucket_name := "commoncrawl"
	file_name := "/crawl-data/CC-MAIN-2016-26/wat.paths.gz"

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(file_name),
		})
	if err != nil {
		fmt.Println("Failed to download file", err)
		return nil
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	// XXX do we really need to close and re-open this??
	inf, _ := os.Open(file.Name())
	_ = inf
	gr, _ := gzip.NewReader(file)
	scanner := bufio.NewScanner(gr)
	for scanner.Scan() {
		indexfiles = append(indexfiles, scanner.Text())
	}

	return indexfiles
}
