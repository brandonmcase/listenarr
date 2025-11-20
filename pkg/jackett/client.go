package jackett

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a Jackett API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Jackett API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query      string
	Category   []int // Category IDs (e.g., 3030 for Books)
	TrackerIDs []string
}

// SearchResult represents a search result from Jackett
type SearchResult struct {
	FirstSeen            time.Time `json:"FirstSeen"`
	Tracker              string    `json:"Tracker"`
	TrackerID            string    `json:"TrackerId"`
	CategoryDesc         string    `json:"CategoryDesc"`
	BlackholeLink        string    `json:"BlackholeLink"`
	Title                string    `json:"Title"`
	Guid                 string    `json:"Guid"`
	Link                 string    `json:"Link"`
	Comments             string    `json:"Comments"`
	PublishDate          time.Time `json:"PublishDate"`
	Category             []int     `json:"Category"`
	Size                 int64     `json:"Size"`
	Files                int       `json:"Files"`
	Grabs                int       `json:"Grabs"`
	Description          string    `json:"Description"`
	RageID               int       `json:"RageID"`
	TVDbID               int       `json:"TVDbId"`
	Imdb                 int       `json:"Imdb"`
	TMDb                 int       `json:"TMDb"`
	Seeders              int       `json:"Seeders"`
	Peers                int       `json:"Peers"`
	MinimumRatio         float64   `json:"MinimumRatio"`
	MinimumSeedTime      int64     `json:"MinimumSeedTime"`
	DownloadVolumeFactor float64   `json:"DownloadVolumeFactor"`
	UploadVolumeFactor   float64   `json:"UploadVolumeFactor"`
	MagnetURI            string    `json:"MagnetUri"`
	InfoHash             string    `json:"InfoHash"`
}

// SearchResponse represents the response from Jackett search
type SearchResponse struct {
	Results  []SearchResult `json:"Results"`
	Indexers []IndexerInfo  `json:"Indexers"`
}

// IndexerInfo represents information about an indexer
type IndexerInfo struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Status  int    `json:"Status"`
	Results int    `json:"Results"`
	Error   string `json:"Error"`
}

// Search performs a search across all configured indexers
func (c *Client) Search(req SearchRequest) (*SearchResponse, error) {
	searchURL := fmt.Sprintf("%s/api/v2.0/indexers/all/results", c.baseURL)

	query := url.Values{}
	query.Set("apikey", c.apiKey)
	query.Set("Query", req.Query)

	// Add category filter (3030 = Books category)
	if len(req.Category) > 0 {
		for _, cat := range req.Category {
			query.Add("Category[]", fmt.Sprintf("%d", cat))
		}
	} else {
		// Default to Books category
		query.Add("Category[]", "3030")
	}

	// Add tracker filter if specified
	if len(req.TrackerIDs) > 0 {
		for _, trackerID := range req.TrackerIDs {
			query.Add("Tracker[]", trackerID)
		}
	}

	searchURL += "?" + query.Encode()

	httpReq, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &searchResp, nil
}

// GetIndexers returns a list of all configured indexers
type Indexer struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Language    string     `json:"language"`
	Encoding    string     `json:"encoding"`
	Categories  []Category `json:"categories"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IndexersResponse struct {
	Indexers []Indexer `json:"indexers"`
}

func (c *Client) GetIndexers() (*IndexersResponse, error) {
	indexersURL := fmt.Sprintf("%s/api/v2.0/indexers/all/results/torznab/api?apikey=%s&t=indexers", c.baseURL, c.apiKey)

	httpReq, err := http.NewRequest("GET", indexersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexers request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get indexers failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Jackett returns XML for indexers, but we'll parse it as JSON if possible
	// For now, return a simple response
	// TODO: Implement proper XML parsing or use a different endpoint
	var indexersResp IndexersResponse
	// This is a placeholder - actual implementation would parse XML
	return &indexersResp, nil
}

// TestConnection tests the connection to Jackett
func (c *Client) TestConnection() error {
	// Try a simple search to test connection
	req := SearchRequest{
		Query: "test",
	}
	_, err := c.Search(req)
	return err
}
