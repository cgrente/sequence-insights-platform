package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ingestRequest struct {
	Values []int64 `json:"values"`
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080/ingest", "Base URL of the ingestion endpoint")
	apiKey := flag.String("api-key", "", "API key / bearer token for authentication")
	cmd := flag.String("cmd", "ingest", "Command to execute: ingest or query")
	valueStr := flag.String("values", "", "Comma-separated list of integer values to ingest")
	id := flag.String("id", "", "ID for querying data")
	flag.Parse()

	hc := &http.Client{Timeout: 10 * time.Second}

	switch *cmd {
	case "ingest":
		values, err := parseValues(*valueStr)
		if err != nil {
			exitErr(err)
		}
		if len(values) == 0 {
			exitErr(fmt.Errorf("no values provided for ingestion"))
		}

		body, _ := json.Marshal(ingestRequest{Values: values})
		req, _ := http.NewRequest(http.MethodPost, strings.TrimRight(*baseURL, "/")+"/v1/sequences/ingest", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if *apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+*apiKey)
		}
		resp, err := hc.Do(req)
		if err != nil {
			exitErr(err)
		}
		defer resp.Body.Close()

		_, _ = io.Copy(os.Stdout, resp.Body)

	case "get":
		if strings.TrimSpace(*id) == "" {
			exitErr(fmt.Errorf("id is required for querying data"))
		}
		req, err := http.NewRequest(http.MethodGet, strings.TrimRight(*baseURL, "/")+"/v1/sequences/"+*id, nil)
		if *apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+*apiKey)
		}
		resp, err := hc.Do(req)
		if err != nil {
			exitErr(err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(os.Stdout, resp.Body)

	default:
		exitErr(fmt.Errorf("unknown command: %s", *cmd))

	}
}

func parseValues(valueStr string) ([]int64, error) {
	valueStr = strings.TrimSpace(valueStr)
	if valueStr == "" {
		return nil, nil
	}
	parts := strings.Split(valueStr, ",")
	out := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		i, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value: %q, %w", part, err)
		}
		out = append(out, i)
	}
	return out, nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
