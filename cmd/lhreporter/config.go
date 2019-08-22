package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type (
	// scoreSet contains members representing each score provided by lighthouse.
	scoreSet struct {
		Performance   int64 `json:"performance"`
		Accessibility int64 `json:"accessibility"`
		BestPractises int64 `json:"bestPractises"`
		SEO           int64 `json:"seo"`
	}

	appConfig struct {
		Remote           string   `json:"remote"`
		MinimumPageScore scoreSet `json:"minimumPageScore"`
		MinimumMeanScore scoreSet `json:"minimumMeanScore"`
		SiteMap          string   `json:"sitemap"`
		CustomPaths      []string `json:"customPaths"`
		StoragePath      string   `json:"storagePath"`
		LighthouseArgs   []string `json:"lighthouseArgs"`
	}
)

func getConfiguration(ctx context.Context) *appConfig {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	c, err := loadConfiguration(ctx, os.Args[1])

	if err != nil {
		log.Fatalf("Usage: lhreporter <config.json>\nError: %s\n", err)
	}

	return c
}

func loadConfiguration(ctx context.Context, loc string) (*appConfig, error) {
	if strings.HasPrefix(loc, "gs://") {
		return loadConfigurationFromBucket(ctx, loc)
	}

	if strings.HasPrefix(loc, "http") {
		return loadConfigurationFromHTTP(ctx, loc)
	}

	return loadConfigurationFromFile(loc)
}

func loadConfigurationFromBucket(ctx context.Context, loc string) (*appConfig, error) {
	u, err := url.Parse(loc)
	if err != nil {
		return nil, err
	}

	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	b := c.Bucket(u.Host)
	r, err := b.Object(u.Path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	rb, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	conf := &appConfig{}
	return conf, json.Unmarshal(rb, conf)
}

func loadConfigurationFromHTTP(ctx context.Context, loc string) (*appConfig, error) {
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rb, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	conf := &appConfig{}
	return conf, json.Unmarshal(rb, conf)
}

func loadConfigurationFromFile(fn string) (*appConfig, error) {
	rb, err := ioutil.ReadFile(fn)

	if err != nil {
		return nil, err
	}

	c := &appConfig{}
	return c, json.Unmarshal(rb, c)
}

// returns true if there is a minimum requirement for any of the scores.
func (s scoreSet) isRequired() bool {
	return s.Performance+s.Accessibility+s.BestPractises+s.SEO > 0
}

// returns true if there is a minimum requirement for any of the scores.
func (s scoreSet) test(result *scoreSet) (string, bool) {
	o := make([]string, 4)
	allPassed := true

	var rT string
	var passed bool

	rT, passed = scoreSetTest("performance", s.Performance, result.Performance)
	o[0] = rT
	allPassed = allPassed && passed

	rT, passed = scoreSetTest("accessibility", s.Accessibility, result.Accessibility)
	o[1] = rT
	allPassed = allPassed && passed

	rT, passed = scoreSetTest("bestPractises", s.BestPractises, result.BestPractises)
	o[2] = rT
	allPassed = allPassed && passed

	rT, passed = scoreSetTest("seo", s.SEO, result.SEO)
	o[3] = rT
	allPassed = allPassed && passed

	return "\t- " + strings.Join(o, "\n\t- ") + "\n", allPassed
}

func scoreSetTest(name string, min, actual int64) (string, bool) {
	if min == 0 {
		return name + " test is not required", true
	}

	if min > actual {
		return fmt.Sprintf(
			"%s requirements failed expected >%d got %d",
			name,
			min,
			actual,
		), false
	}

	return fmt.Sprintf(
		"%s requirements passed expected >%d got %d",
		name,
		min,
		actual,
	), true
}
