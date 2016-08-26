package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danohuiginn/grepurl/grepurl"
	"github.com/urfave/cli"
)

func NotImplemented(c *cli.Context) error {
	fmt.Println("Not implemented yet")

	return nil
}

func Sample(cli *cli.Context) error {
	files := []string{"/tmp/urlsample.txt"}
	urlstore, trigrams := grepurl.RunImport(files)
	fmt.Println("import complete")
	query := ".*.gov"
	ch := make(chan string)
	go grepurl.RunQuery(query, trigrams, urlstore, ch)
	for result := range ch {
		fmt.Println(result)
	}
	return nil
}

func Download(cli *cli.Context) error {
	/*cmd := exec.Command("/home/dan/.virtualenvs/yl/bin/aws", "s3", "cp", "s3://commoncrawl/crawl-data/CC-MAIN-2016-22/segments/1464049270134.8/wat/CC-MAIN-20160524002110-00000-ip-10-185-217-139.ec2.internal.warc.wat.gz", "/tmp/mydl.gz")

	output, err := cmd.CombinedOutput()
	log.Println(strings.Join(cmd.Args, " "))
	log.Println(err)
	log.Println(string(output))
	return nil
	fn := "crawl-data/CC-MAIN-2016-22/segments/1464049270134.8/wat/CC-MAIN-20160524002110-00000-ip-10-185-217-139.ec2.internal.warc.wat.gz"
	fn = "/tmp/unzip"*/
	log.Println("starting upload")
	indices := grepurl.ListIndexFiles("2016-05")
	linecount := grepurl.UploadURLs(indices[0], true)
	log.Println("uploaded " + strconv.Itoa(linecount) + " lines.")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "grepurl"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir",
			Value: "/tmp/grepurl/",
			Usage: "Directory for grepurl to keep its data. Should have plenty of space, and ideally be accessible from multiple machines",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "search",
			Usage:  "look for an url",
			Action: NotImplemented,
		},
		{
			Name:   "sample",
			Usage:  "give it a spin",
			Action: Sample,
		},
		{
			Name:   "download",
			Usage:  "download some crawl archives",
			Action: Download,
		},
	}
	app.Run(os.Args)

}
