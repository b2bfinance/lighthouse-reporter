package main

import (
	"context"
	"os/exec"

	"github.com/tidwall/gjson"
)

func getLighthouseScoreForEndpoint(ctx context.Context, endpoint string) (*scoreSet, error) {
	cmd := exec.CommandContext(ctx, "lighthouse", endpoint, "--quiet", "--output", "json")

	rb, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return &scoreSet{
		Performance:   convertPercentageToInt(gjson.GetBytes(rb, "categories.performance.score").Float()),
		Accessibility: convertPercentageToInt(gjson.GetBytes(rb, "categories.accessibility.score").Float()),
		BestPractises: convertPercentageToInt(gjson.GetBytes(rb, "categories.best-practices.score").Float()),
		SEO:           convertPercentageToInt(gjson.GetBytes(rb, "categories.seo.score").Float()),
	}, nil
}

func convertPercentageToInt(in float64) int64 {
	return int64(in * 100)
}
