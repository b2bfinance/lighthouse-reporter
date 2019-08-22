package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	started time.Time
)

func main() {
	started := time.Now()

	ctx := context.Background()

	if len(os.Args) < 2 {
		log.Fatalln("Usage: lhreporter <config.json> [reference]")
	}

	c := getConfiguration(ctx)

	var reference string
	if len(os.Args) > 2 {
		reference = os.Args[2]
	} else {
		reference = fmt.Sprintf("%d", started.UnixNano())
	}

	checkLighthouseExists()
	checkRemoteIsResponsive(ctx, c.Remote)

	remote := strings.TrimRight(c.Remote, "/")

	endpoints := make([]string, len(c.CustomPaths))
	for i, ep := range c.CustomPaths {
		endpoints[i] = remote + "/" + strings.TrimLeft(ep, "/")
	}

	if c.SiteMap != "" {
		e, err := extractURLsFromSitemap(ctx, remote, c.SiteMap)
		if err != nil {
			log.Fatalln("unable to extract URLs from sitemap: ", err)
		}
		endpoints = append(endpoints, e...)
	}

	epoffset := len(remote)
	results := make(endpointResults, 0, len(endpoints))
	for _, ep := range endpoints {
		score, err := getLighthouseScoreForEndpoint(ctx, c.LighthouseArgs, ep)
		if err != nil {
			log.Printf("unable to get result for '%s' error: %s", ep, err)
			continue
		}
		results = append(results, endpointResult{
			Path:  ep[epoffset:],
			Score: score,
		})
	}

	if err := outputResults(ctx, reference, c.StoragePath, results); err != nil {
		log.Fatalln("unable to store results: ", err)
	}

	failed := false

	if c.MinimumPageScore.isRequired() {
		log.Println("----------------------------------")
		log.Println("-- Minimum page score results")
		log.Println("----------------------------------")

		for _, r := range results {
			res, passed := c.MinimumPageScore.test(r.Score)
			log.Print("Endpoint: ", r.Path)
			log.Print(res)
			log.Println("")
			if !passed {
				failed = true
			}
		}
	}

	if c.MinimumMeanScore.isRequired() {
		log.Println("----------------------------------")
		log.Println("-- Minimum mean score results")
		log.Println("----------------------------------")

		res, passed := c.MinimumMeanScore.test(results.Mean())
		log.Print(res)
		if !passed {
			failed = true
		}
	}

	if failed {
		log.Fatalln("Some score requirements failed, please see the output above for more information.")
	}
}

func checkLighthouseExists() {
	cmd := exec.Command("lighthouse", "--version")
	if err := cmd.Run(); err != nil {
		log.Println("Unable to find NPM package lighthouse make sure it is in your $PATH")
		log.Fatalln("Error:", err)
	}
}

func checkRemoteIsResponsive(ctx context.Context, remote string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req, err := http.NewRequest("GET", remote, nil)
	if err != nil {
		log.Fatalf("unable to check remote, maybe a malformed remote? Error: %s\n", err)
	}

	_, err = http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatalf("unable to check remote, maybe it is down? Error: %s\n", err)
	}
}
