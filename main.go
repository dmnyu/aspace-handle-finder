package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
)

var (
	config      string
	environment string
	httpHandle  = "http://hdl.handle.net/2333.1/"
	httpsHandle = "https://hdl.handle.net/2333.1/"
	httpCount   = 0
	httpsCount  = 0
)

func init() {
	flag.StringVar(&config, "config", "", "")
	flag.StringVar(&environment, "environment", "", "")
}

func main() {
	flag.Parse()
	asclient, err := aspace.NewClient(config, environment, 20)
	if err != nil {
		panic(err)
	}

	httpFile, _ := os.Create("http-handles.txt")
	defer httpFile.Close()
	httpWriter := bufio.NewWriter(httpFile)

	httpsFile, _ := os.Create("https-handles.txt")
	defer httpsFile.Close()
	httpsWriter := bufio.NewWriter(httpsFile)

	for _, repoID := range []int{2, 3, 6} {
		initialRequest, err := asclient.Search(repoID, "digital_object", "*", 1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Repository %d\n", repoID)
		for pageNum := initialRequest.FirstPage; pageNum <= initialRequest.LastPage; pageNum++ {
			fmt.Printf("page %d of %d\n", pageNum, initialRequest.LastPage)
			request, err := asclient.Search(repoID, "digital_object", "*", pageNum)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, resource := range request.Results {
				do := aspace.DigitalObject{}
				do_bytes := []byte(fmt.Sprint(resource["json"]))
				err := json.Unmarshal(do_bytes, &do)
				if err != nil {
					log.Println(err)
				}
				for _, fv := range do.FileVersions {
					if strings.Contains(fv.FileURI, httpHandle) {
						httpCount++
						httpWriter.WriteString(fmt.Sprintf("%s %s\n", do.URI, fv.FileURI))
					} else if strings.Contains(fv.FileURI, httpsHandle) {
						httpsCount++
						httpsWriter.WriteString(fmt.Sprintf("%s %s\n", do.URI, fv.FileURI))
					}
				}
			}
		}
	}
	fmt.Printf("http: %d\thttps: %d\n", httpCount, httpsCount)
}
