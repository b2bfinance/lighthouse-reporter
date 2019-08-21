package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEndpointFromSitemapURL(t *testing.T) {
	tcs := []struct {
		Base     string
		In       string
		Expected string
	}{
		{
			Base:     "https://b2bfinance.com",
			In:       "https://b2bfinance.com/page-a",
			Expected: "https://b2bfinance.com/page-a",
		},
		{
			Base:     "https://b2bfinance.com",
			In:       "https://example.com/page-a",
			Expected: "https://b2bfinance.com/page-a",
		},
		{
			Base:     "https://b2bfinance.com",
			In:       "/page-a",
			Expected: "https://b2bfinance.com/page-a",
		},
		{
			Base:     "https://b2bfinance.com",
			In:       "page-a",
			Expected: "https://b2bfinance.com/page-a",
		},
		{
			Base:     "https://b2bfinance.com",
			In:       "/page-a?b=1",
			Expected: "https://b2bfinance.com/page-a?b=1",
		},
		{
			Base:     "https://b2bfinance.com",
			In:       "https://b2bfinance.com/page-a?b=1",
			Expected: "https://b2bfinance.com/page-a?b=1",
		},
	}

	for i, tc := range tcs {
		actual := buildEndpointFromSitemapURL(tc.Base, tc.In)

		assert.Equal(t, tc.Expected, actual, fmt.Sprintf("Test case %d", i))
	}
}
