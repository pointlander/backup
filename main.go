// Copyright 2026 The Backup Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Repository represents the minimal structure of a GitHub repository payload
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

func main() {
	page := 1
	for {
		url := fmt.Sprintf("https://api.github.com/users/pointlander/repos?per_page=100&page=%d", page)
		client := &http.Client{
			Timeout: 15 * time.Second,
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("User-Agent", "Go-GitHub-Client")
		req.Header.Set("Accept", "application/vnd.github+json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			panic(fmt.Errorf("API returned error status %d: %s\n", resp.StatusCode, string(body)))
		}
		var repos []Repository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			panic(fmt.Errorf("Error decoding JSON payload: %v\n", err))
		}

		if len(repos) == 0 {
			break
		}

		for i, repo := range repos {
			fmt.Printf("%d. %s\n   URL: %s\n   Desc: %s\n\n", (page-1)*100+i+1, repo.Name, repo.HTMLURL, repo.Description)
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/pointlander/%s.git", repo.Name))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				panic(fmt.Errorf("Error cloning repository: %v\n", err))
			}
		}
		page++
	}
}
