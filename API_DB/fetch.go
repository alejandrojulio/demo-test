package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func fetchAndStoreData(db *DB, baseURL, page string) error {
	url := baseURL
	if page != "" {
		url += "?next_page=" + page
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error HTTP: %w", err)
	}

	token := os.Getenv("API_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error GET: %w", err)
	}
	defer resp.Body.Close()

	var data ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error decode: %w", err)
	}

	for _, item := range data.Items {
		if err := db.InsertItem(item); err != nil {
			return fmt.Errorf("error insert: %w", err)
		}
	}

	if data.NextPage != "" {
		return fetchAndStoreData(db, baseURL, data.NextPage)
	}

	return nil
}
