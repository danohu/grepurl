package main

import (
	"fmt"
	"os"

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
	grepurl.ListIndexFiles("2016-05")
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
