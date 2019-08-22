package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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

func loadFromFile(fn string) (*appConfig, error) {
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
