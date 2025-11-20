package jackett

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:9117", "test-api-key")
	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:9117", client.baseURL)
	assert.Equal(t, "test-api-key", client.apiKey)
}

func TestClient_Search(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2.0/indexers/all/results" {
			// Check API key
			apiKey := r.URL.Query().Get("apikey")
			if apiKey != "test-api-key" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Return mock search results
			results := SearchResponse{
				Results: []SearchResult{
					{
						Title:       "Test Audiobook",
						Tracker:     "TestTracker",
						TrackerID:   "test-tracker",
						Size:        1000000000,
						Seeders:     10,
						Peers:       15,
						MagnetURI:   "magnet:?xt=urn:btih:test123",
						InfoHash:    "test123",
						PublishDate: time.Now(),
					},
				},
				Indexers: []IndexerInfo{
					{
						ID:      "test-tracker",
						Name:    "TestTracker",
						Status:  1,
						Results: 1,
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")

	req := SearchRequest{
		Query: "test audiobook",
	}

	resp, err := client.Search(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 1)
	assert.Equal(t, "Test Audiobook", resp.Results[0].Title)
}

func TestClient_Search_NoResults(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2.0/indexers/all/results" {
			results := SearchResponse{
				Results:  []SearchResult{},
				Indexers: []IndexerInfo{},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")

	req := SearchRequest{
		Query: "nonexistent book",
	}

	resp, err := client.Search(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 0)
}

func TestClient_TestConnection(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2.0/indexers/all/results" {
			results := SearchResponse{
				Results:  []SearchResult{},
				Indexers: []IndexerInfo{},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	err := client.TestConnection()
	assert.NoError(t, err)
}
