package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

func outputResults(ctx context.Context, reference, loc string, r endpointResults) error {
	if strings.HasPrefix(loc, "gs://") {
		return outputResultsToBucket(ctx, reference, loc, r)
	}

	if strings.HasPrefix(loc, "http") {
		return outputResultsToHTTP(ctx, reference, loc, r)
	}

	return outputResultsToFile(loc, reference, r)
}

func outputResultsToBucket(ctx context.Context, reference, loc string, r endpointResults) error {
	u, err := url.Parse(loc)
	if err != nil {
		return err
	}

	c, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	path := strings.Trim(u.Path, "/")
	if u.Path != "" {
		path += "/"
	} 

	b := c.Bucket(u.Host)
	w := b.Object(path + reference + ".json").NewWriter(ctx)

	rb, _ := json.Marshal(struct {
		Reference   string          `json:"reference"`
		Results     endpointResults `json:"results"`
		MeanResults *scoreSet       `json:"meanResults"`
		Started     int64           `json:"started"`
		Finished    int64           `json:"finished"`
	}{
		Reference:   reference,
		Results:     r,
		MeanResults: r.Mean(),
		Started:     started.Unix(),
		Finished:    time.Now().Unix(),
	})

	if _, err := io.Copy(w, bytes.NewReader(rb)); err != nil {
		return err
	}

	return w.Close()
}

func outputResultsToHTTP(ctx context.Context, reference, loc string, r endpointResults) error {
	rb, _ := json.Marshal(struct {
		Reference   string          `json:"reference"`
		Results     endpointResults `json:"results"`
		MeanResults *scoreSet       `json:"meanResults"`
		Started     int64           `json:"started"`
		Finished    int64           `json:"finished"`
	}{
		Reference:   reference,
		Results:     r,
		MeanResults: r.Mean(),
		Started:     started.Unix(),
		Finished:    time.Now().Unix(),
	})

	req, err := http.NewRequest("POST", loc, bytes.NewReader(rb))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"error uploading results expected status code %d got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	return nil
}

func outputResultsToFile(reference, loc string, r endpointResults) error {
	rb, _ := json.Marshal(struct {
		Reference   string          `json:"reference"`
		Results     endpointResults `json:"results"`
		MeanResults *scoreSet       `json:"meanResults"`
		Started     int64           `json:"started"`
		Finished    int64           `json:"finished"`
	}{
		Reference:   reference,
		Results:     r,
		MeanResults: r.Mean(),
		Started:     started.Unix(),
		Finished:    time.Now().Unix(),
	})

	return ioutil.WriteFile(loc, rb, os.FileMode(0666))
}
