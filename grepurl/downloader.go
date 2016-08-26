package grepurl

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
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

var TESTING = false

const bucket_name = "commoncrawl"
const s3_upload_bucket = "alephdev.openoil.net"
const s3_upload_prefix = "urls/"
const s3_upload_region = "eu-west-1"

const redis_host = "ohuiginn.net"
const redis_port = "90995"
const redis_password = "Sheu9dooUn1eech0"
const redis_db = 9
const redis_name = "grepurl"

//XXX month should look like "2016-06"
// beware tat
// -- just getting the latest month for onw
func ListIndexFiles(month string) []string {

	indexfiles := []string{}

	monthindex := MonthIndexes[month]
	file_name := "/crawl-data/CC-MAIN-" + monthindex + "/wat.paths.gz"
	fmt.Println(file_name)

	scanner, file := s3Gzip(bucket_name, file_name)
	defer file.Close()
	defer os.Remove(file.Name())
	for scanner.Scan() {
		indexfiles = append(indexfiles, scanner.Text())
	}
	return indexfiles
}

func PushIndexFiles() {
}

/* This is assuming basic common_crawl stuff

 */
func URLsFromWAT(path string, ch chan string) {
	scanner, file := s3Gzip(bucket_name, path)
	defer file.Close()
	defer os.Remove(file.Name())

	// for handling of long lines
	buf := make([]byte, 0, 64*1024*1000)
	scanner.Buffer(buf, 1024*1024*100)

	lines := 0
	for scanner.Scan() {
		line := scanner.Text()
		lines += 1
		if lines%1000 == 0 {
			log.Println("found " + strconv.Itoa(lines) + " lines")
		}

		if strings.HasPrefix(line, "WARC-Target-URI: ") {
			url := strings.SplitN(line, ": ", 2)
			ch <- url[1]
		}

	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	close(ch)
}

func GenerateS3Path(origpath string) string {
	// week of year
	rx, _ := regexp.Compile(`CC-MAIN-([\d\-]+)`)
	match := rx.FindStringSubmatch(origpath)
	week := match[1]

	// filename, based on date ad IP of scraper
	rx, _ = regexp.Compile(`/([^/]+)ec2.internal.warc.wat.gz$`)
	fn := rx.FindStringSubmatch(origpath)[1]
	filepath := s3_upload_prefix + week + "/" + fn + "urls.gz"
	return filepath
}

func S3PathExists(bucket string, path string, region string) bool {
	sess := session.New(&aws.Config{Region: aws.String(region)})
	svc := s3.New(sess)
	params := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket), // Required
		MaxKeys: aws.Int64(1),
		Prefix:  aws.String(path),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Println("s3 error: ")
		fmt.Println(err)
		fmt.Println(bucket)
	}
	return len(resp.Contents) > 0
}

/*
Return the number of URLs uploaded,
or -1 if we didn't upload anything
(e.g. because we avoided replacing an existing file)
*/
func UploadURLs(dl_path string, clobber bool) int {
	log.Println("uploadurls happening")
	upload_path := ""
	if TESTING {
		upload_path = "urls/TEST_TEST_TEST"
	} else {
		upload_path = GenerateS3Path(dl_path)
	}
	if !clobber && S3PathExists(s3_upload_bucket, upload_path, s3_upload_region) {
		return -1
	}
	log.Println("upload urls slowly chugging")
	urls := make(chan string)
	go URLsFromWAT(dl_path, urls)
	log.Println("started urlsfromwat")
	urlfile, _ := ioutil.TempFile("", "uploadurls_")
	log.Println("urlfile created as " + urlfile.Name())

	defer urlfile.Close()
	defer os.Remove(urlfile.Name())
	lines := 0

	log.Println("iterating through urls")
	for urlstr := range urls {
		lines += 1
		io.WriteString(urlfile, urlstr+"\n")
	}
	log.Println("iterated through urls")
	urlfile.Sync()

	sess := session.New(&aws.Config{Region: aws.String(s3_upload_region)})
	uploader := s3manager.NewUploader(sess)

	//uploadFile, _ := os.Open(urlfile.Name())
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(s3_upload_bucket),
		Key:    aws.String(upload_path),
		Body:   urlfile,
	}
	_, err := uploader.Upload(upParams)
	if err != nil {
		log.Println("upload error")
		log.Println(err)
		return -1
	}
	return lines
}

// Take a gzipped s3 file, and return a scanner over it
// Internally the file is downloaded to a temporary location,
// unzipped and then removed
func s3Gzip(bucket_name string, file_name string) (*bufio.Scanner, *os.File) {
	file, _ := ioutil.TempFile("", "s3gzip_")
	//defer file.Close()
	//defer os.Remove(file.Name())
	log.Println("s3gzip " + bucket_name + " " + file_name)

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(file_name),
		})
	log.Println("s3gzip done" + bucket_name + " " + file_name)
	gr, _ := gzip.NewReader(file)
	scanner := bufio.NewScanner(gr)
	return scanner, file
}
