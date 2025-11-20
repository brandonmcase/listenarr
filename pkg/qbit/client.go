package qbit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a qBittorrent API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	sid        string // Session ID
}

// NewClient creates a new qBittorrent API client
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: username,
		password: password,
	}
}

// Login authenticates with qBittorrent and stores the session ID
func (c *Client) Login() error {
	loginURL := fmt.Sprintf("%s/api/v2/auth/login", c.baseURL)

	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check response body for "Ok." or "Fails."
	responseText := strings.TrimSpace(string(body))
	if responseText != "Ok." {
		return fmt.Errorf("login failed: %s", responseText)
	}

	// Extract session ID from cookies
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			c.sid = cookie.Value
			return nil
		}
	}

	// If no SID cookie, try to get it from Set-Cookie header
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie != "" {
		parts := strings.Split(setCookie, ";")
		for _, part := range parts {
			if strings.HasPrefix(strings.TrimSpace(part), "SID=") {
				c.sid = strings.TrimPrefix(strings.TrimSpace(part), "SID=")
				return nil
			}
		}
	}

	return fmt.Errorf("no session ID received from qBittorrent")
}

// Logout logs out from qBittorrent
func (c *Client) Logout() error {
	logoutURL := fmt.Sprintf("%s/api/v2/auth/logout", c.baseURL)

	req, err := http.NewRequest("POST", logoutURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.sid = ""
	return nil
}

// setAuthHeader sets the authentication cookie on the request
func (c *Client) setAuthHeader(req *http.Request) {
	if c.sid != "" {
		req.AddCookie(&http.Cookie{
			Name:  "SID",
			Value: c.sid,
		})
	}
}

// AddTorrent adds a torrent to qBittorrent
func (c *Client) AddTorrent(torrentURL string, options *AddTorrentOptions) error {
	addURL := fmt.Sprintf("%s/api/v2/torrents/add", c.baseURL)

	data := url.Values{}
	data.Set("urls", torrentURL)

	if options != nil {
		if options.Category != "" {
			data.Set("category", options.Category)
		}
		if options.SavePath != "" {
			data.Set("savepath", options.SavePath)
		}
		if options.Paused {
			data.Set("paused", "true")
		}
		if options.RootFolder {
			data.Set("root_folder", "true")
		}
		if options.Rename != "" {
			data.Set("rename", options.Rename)
		}
		if options.UploadLimit > 0 {
			data.Set("upLimit", fmt.Sprintf("%d", options.UploadLimit))
		}
		if options.DownloadLimit > 0 {
			data.Set("dlLimit", fmt.Sprintf("%d", options.DownloadLimit))
		}
		if options.SequentialDownload {
			data.Set("sequentialDownload", "true")
		}
		if options.FirstLastPiecePriority {
			data.Set("firstLastPiecePrio", "true")
		}
		if options.SkipChecking {
			data.Set("skip_checking", "true")
		}
		if options.ContentLayout != "" {
			data.Set("contentLayout", options.ContentLayout)
		}
		if options.AutoTMM {
			data.Set("autoTMM", "true")
		}
	}

	req, err := http.NewRequest("POST", addURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create add torrent request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add torrent failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// AddTorrentOptions represents options for adding a torrent
type AddTorrentOptions struct {
	Category               string
	SavePath               string
	Paused                 bool
	RootFolder             bool
	Rename                 string
	UploadLimit            int64 // bytes per second
	DownloadLimit          int64 // bytes per second
	SequentialDownload     bool
	FirstLastPiecePriority bool
	SkipChecking           bool
	ContentLayout          string // "Original", "Subfolder", "NoSubfolder"
	AutoTMM                bool   // Automatic Torrent Management
}

// TorrentInfo represents information about a torrent
type TorrentInfo struct {
	Hash          string  `json:"hash"`
	Name          string  `json:"name"`
	Size          int64   `json:"size"`
	Progress      float64 `json:"progress"` // 0-1
	State         string  `json:"state"`
	Downloaded    int64   `json:"downloaded"`
	Uploaded      int64   `json:"uploaded"`
	DownloadSpeed int64   `json:"dlspeed"` // bytes per second
	UploadSpeed   int64   `json:"upspeed"` // bytes per second
	ETA           int64   `json:"eta"`     // seconds
	Category      string  `json:"category"`
	SavePath      string  `json:"save_path"`
	ContentPath   string  `json:"content_path"`
	AddedOn       int64   `json:"added_on"`
	CompletionOn  int64   `json:"completion_on"`
	Tracker       string  `json:"tracker"`
	Seeds         int     `json:"num_seeds"`
	Leechers      int     `json:"num_leechs"`
	Ratio         float64 `json:"ratio"`
}

// GetTorrentList returns a list of all torrents
func (c *Client) GetTorrentList(filters *TorrentFilters) ([]TorrentInfo, error) {
	listURL := fmt.Sprintf("%s/api/v2/torrents/info", c.baseURL)

	if filters != nil {
		query := url.Values{}
		if filters.Category != "" {
			query.Set("category", filters.Category)
		}
		if filters.Filter != "" {
			query.Set("filter", filters.Filter)
		}
		if filters.Sort != "" {
			query.Set("sort", filters.Sort)
		}
		if filters.Reverse {
			query.Set("reverse", "true")
		}
		if filters.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", filters.Limit))
		}
		if filters.Offset > 0 {
			query.Set("offset", fmt.Sprintf("%d", filters.Offset))
		}
		if len(query) > 0 {
			listURL += "?" + query.Encode()
		}
	}

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create torrent list request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get torrent list failed with status %d: %s", resp.StatusCode, string(body))
	}

	var torrents []TorrentInfo
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, fmt.Errorf("failed to decode torrent list: %w", err)
	}

	return torrents, nil
}

// TorrentFilters represents filters for getting torrent list
type TorrentFilters struct {
	Category string
	Filter   string // "all", "downloading", "completed", "paused", "active", "inactive", "resumed", "stalled", "stalled_uploading", "stalled_downloading"
	Sort     string // "name", "size", "progress", "dlspeed", "upspeed", "priority", "added_on", "completion_on", "tracker", "state"
	Reverse  bool
	Limit    int
	Offset   int
}

// GetTorrentInfo returns information about a specific torrent by hash
func (c *Client) GetTorrentInfo(hash string) (*TorrentInfo, error) {
	torrents, err := c.GetTorrentList(nil)
	if err != nil {
		return nil, err
	}

	for i := range torrents {
		if strings.EqualFold(torrents[i].Hash, hash) {
			return &torrents[i], nil
		}
	}

	return nil, fmt.Errorf("torrent with hash %s not found", hash)
}

// DeleteTorrent deletes a torrent from qBittorrent
func (c *Client) DeleteTorrent(hashes []string, deleteFiles bool) error {
	deleteURL := fmt.Sprintf("%s/api/v2/torrents/delete", c.baseURL)

	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	if deleteFiles {
		data.Set("deleteFiles", "true")
	}

	req, err := http.NewRequest("POST", deleteURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete torrent failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// PauseTorrent pauses one or more torrents
func (c *Client) PauseTorrent(hashes []string) error {
	return c.torrentAction("pause", hashes)
}

// ResumeTorrent resumes one or more torrents
func (c *Client) ResumeTorrent(hashes []string) error {
	return c.torrentAction("resume", hashes)
}

// torrentAction performs a generic torrent action
func (c *Client) torrentAction(action string, hashes []string) error {
	actionURL := fmt.Sprintf("%s/api/v2/torrents/%s", c.baseURL, action)

	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", actionURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", action, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to %s torrent: %w", action, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s torrent failed with status %d: %s", action, resp.StatusCode, string(body))
	}

	return nil
}

// GetTorrentProperties returns detailed properties of a torrent
type TorrentProperties struct {
	SavePath           string  `json:"save_path"`
	CreationDate       int64   `json:"creation_date"`
	PieceSize          int64   `json:"piece_size"`
	Comment            string  `json:"comment"`
	TotalWasted        int64   `json:"total_wasted"`
	TotalUploaded      int64   `json:"total_uploaded"`
	TotalDownloaded    int64   `json:"total_downloaded"`
	UpLimit            int64   `json:"up_limit"`
	DlLimit            int64   `json:"dl_limit"`
	TimeElapsed        int64   `json:"time_elapsed"`
	SeedingTime        int64   `json:"seeding_time"`
	NbConnections      int     `json:"nb_connections"`
	NbConnectionsLimit int     `json:"nb_connections_limit"`
	ShareRatio         float64 `json:"share_ratio"`
	AdditionDate       int64   `json:"addition_date"`
	CompletionDate     int64   `json:"completion_date"`
	CreatedBy          string  `json:"created_by"`
	DlSpeedAvg         int64   `json:"dl_speed_avg"`
	UpSpeedAvg         int64   `json:"up_speed_avg"`
	Eta                int64   `json:"eta"`
	LastSeen           int64   `json:"last_seen"`
	Peers              int     `json:"peers"`
	Seeds              int     `json:"seeds"`
}

func (c *Client) GetTorrentProperties(hash string) (*TorrentProperties, error) {
	propsURL := fmt.Sprintf("%s/api/v2/torrents/properties?hash=%s", c.baseURL, hash)

	req, err := http.NewRequest("GET", propsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create properties request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent properties: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get torrent properties failed with status %d: %s", resp.StatusCode, string(body))
	}

	var props TorrentProperties
	if err := json.NewDecoder(resp.Body).Decode(&props); err != nil {
		return nil, fmt.Errorf("failed to decode torrent properties: %w", err)
	}

	return &props, nil
}

// GetGlobalTransferInfo returns global transfer information
type GlobalTransferInfo struct {
	DlInfoSpeed      int64  `json:"dl_info_speed"`     // Global download speed (bytes/s)
	DlInfoData       int64  `json:"dl_info_data"`      // Data downloaded this session (bytes)
	UpInfoSpeed      int64  `json:"up_info_speed"`     // Global upload speed (bytes/s)
	UpInfoData       int64  `json:"up_info_data"`      // Data uploaded this session (bytes)
	DlRateLimit      int64  `json:"dl_rate_limit"`     // Download rate limit (bytes/s)
	UpRateLimit      int64  `json:"up_rate_limit"`     // Upload rate limit (bytes/s)
	ConnectionStatus string `json:"connection_status"` // Connection status
}

func (c *Client) GetGlobalTransferInfo() (*GlobalTransferInfo, error) {
	infoURL := fmt.Sprintf("%s/api/v2/transfer/info", c.baseURL)

	req, err := http.NewRequest("GET", infoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer info request: %w", err)
	}

	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get transfer info failed with status %d: %s", resp.StatusCode, string(body))
	}

	var info GlobalTransferInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode transfer info: %w", err)
	}

	return &info, nil
}
