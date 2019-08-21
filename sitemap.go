// Sitemap will parse and provide the URLs for a sitemap that is
// compatible with the following protocol specification:
// https://www.sitemaps.org/protocol.html
package main

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	sitemapBase struct {
		URLs []*sitemapURL `xml:"url"`
	}

	sitemapURL struct {
		Loc string `xml:"loc"`
	}
)

func extractURLsFromSitemap(ctx context.Context, base, loc string) ([]string, error) {
	var sm *sitemapBase
	if strings.HasPrefix(loc, "http") { // covers http and https
		s, err := sitemapFromURL(ctx, loc)
		if err != nil {
			return nil, err
		}
		sm = s
	} else {
		s, err := sitemapFromFile(loc)
		if err != nil {
			return nil, err
		}
		sm = s
	}

	res := make([]string, len(sm.URLs))
	for i, u := range sm.URLs {
		res[i] = buildEndpointFromSitemapURL(base, u.Loc)
	}

	return res, nil
}

func buildEndpointFromSitemapURL(base, loc string) string {
	if strings.HasPrefix(loc, "http") {
		u, err := url.Parse(loc)
		if err == nil {
			loc = u.Path
			if len(u.RawQuery) > 0 {
				loc += "?" + u.Query().Encode()
			}
		}
	}

	return base + "/" + strings.TrimLeft(loc, "/")
}

func sitemapFromFile(loc string) (*sitemapBase, error) {
	rb, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}

	return sitemapFromBytes(rb)
}

func sitemapFromURL(ctx context.Context, loc string) (*sitemapBase, error) {
	req, err := http.NewRequest("GET", loc, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return sitemapFromBytes(rb)
}

func sitemapFromBytes(rb []byte) (*sitemapBase, error) {
	sm := &sitemapBase{}
	return sm, xml.Unmarshal(rb, sm)
}
