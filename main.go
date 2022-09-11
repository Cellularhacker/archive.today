package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	optionDownload = flag.String("download", "", "Also wait and download the zip archive.")
	optionBundle   = flag.String("bundle", "", "Merge all the zip archives in a tarball with a HTML index.")
)

func init() {
	flag.Parse()
}

func main() {
	if optionDownload == nil && optionBundle == nil {
		log.Fatal("You must specify a single option.")
	}

	if optionBundle != nil {
		archiveUrl := commit(*optionBundle)
		download(archiveUrl)
	}

	if optionDownload != nil {
		download(*optionDownload)
	}
}

func commit(target string) string {
	p := url.Values{}
	p.Set("url", target)

	resp, err := http.PostForm("https://archive.today/submit/", p)
	if err != nil {
		log.Fatal("Error doing a POST:", err)
	}
	resp.Body.Close()

	h := resp.Header.Get("Refresh")
	if h[:6] != "0;url=" {
		log.Fatal("Malformed answer while committing.")
	}

	return h[6:]
}

var loadingGif = []byte("https://archive.today/loading.gif")

func fetchZip(archiveURL string) (io.Reader, error) {
	zipURL := archiveURL + ".zip"

	for {
		body, err := get(archiveURL)
		if err != nil {
			log.Fatal("Error while checking", zipURL, "-", err)
		}
		if bytes.Index(body, loadingGif) > -1 {
			time.Sleep(1 * time.Second)
			continue
		}

		body, err = get(zipURL)
		if err != nil {
			log.Fatal("Error while downloading", zipURL, "-", err)
		}

		return bytes.NewBuffer(body), nil
	}
}

func download(archiveURL string) {
	tokens := strings.Split(archiveURL, "/")
	fileName := tokens[len(tokens)-1] + ".zip"

	output, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Error while creating", fileName, "-", err)
	}
	defer output.Close()

	respBody, err := fetchZip(archiveURL)
	if err != nil {
		log.Errorf("archive.today failed to create an archive - check %s\n", archiveURL)
		return
	}

	_, err = io.Copy(output, respBody)
	if err != nil {
		log.Fatal("Error while downloading", archiveURL, "-", err)
	}
}

func get(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("archive.today error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
