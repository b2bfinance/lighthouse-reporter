package main

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/tidwall/gjson"
)

func getLighthouseScoreForEndpoint(ctx context.Context, lhArgs []string, endpoint string) (*scoreSet, error) {
	args := []string{
		endpoint,
		"--quiet",
		"--output", "json",
	}
	args = append(args, lhArgs...)

	if v := os.Getenv("LIGHTHOUSE_ARGS"); v != "" {
		args = append(args, strings.Split(v, " ")...)
	}

	cmd := exec.CommandContext(ctx, "lighthouse", args...)

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
