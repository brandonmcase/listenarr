package qbit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8080", "admin", "adminadmin")
	assert.NotNil(t, client)
	assert.Equal(t, "http://localhost:8080", client.baseURL)
	assert.Equal(t, "admin", client.username)
	assert.Equal(t, "adminadmin", client.password)
}

func TestClient_Login(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/auth/login" {
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Check credentials
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "admin" && password == "adminadmin" {
				http.SetCookie(w, &http.Cookie{
					Name:  "SID",
					Value: "test-session-id",
				})
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Ok."))
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Fails."))
			}
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "adminadmin")
	err := client.Login()

	assert.NoError(t, err)
	assert.Equal(t, "test-session-id", client.sid)
}

func TestClient_Login_Failure(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/auth/login" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Fails."))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "wrongpassword")
	err := client.Login()

	assert.Error(t, err)
	assert.Empty(t, client.sid)
}

func TestClient_AddTorrent(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/torrents/add" {
			// Check for SID cookie
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value != "test-session-id" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			r.ParseForm()
			urls := r.FormValue("urls")
			if urls == "magnet:?xt=urn:btih:test" {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else if r.URL.Path == "/api/v2/auth/login" {
			http.SetCookie(w, &http.Cookie{
				Name:  "SID",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "adminadmin")
	err := client.Login()
	assert.NoError(t, err)

	err = client.AddTorrent("magnet:?xt=urn:btih:test", nil)
	assert.NoError(t, err)
}

func TestClient_GetTorrentList(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/torrents/info" {
			// Check for SID cookie
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value != "test-session-id" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"hash":"abc123","name":"Test Torrent","size":1000000,"progress":0.5,"state":"downloading","downloaded":500000,"uploaded":0,"dlspeed":1000000,"upspeed":0,"eta":60,"category":"","save_path":"/downloads","content_path":"/downloads/Test Torrent","added_on":1234567890,"completion_on":0,"tracker":"http://tracker.example.com","num_seeds":10,"num_leechs":5,"ratio":0.0}]`))
		} else if r.URL.Path == "/api/v2/auth/login" {
			http.SetCookie(w, &http.Cookie{
				Name:  "SID",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "adminadmin")
	err := client.Login()
	assert.NoError(t, err)

	torrents, err := client.GetTorrentList(nil)
	assert.NoError(t, err)
	assert.Len(t, torrents, 1)
	assert.Equal(t, "abc123", torrents[0].Hash)
	assert.Equal(t, "Test Torrent", torrents[0].Name)
	assert.Equal(t, 0.5, torrents[0].Progress)
}

func TestClient_GetTorrentInfo(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/torrents/info" {
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value != "test-session-id" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"hash":"abc123","name":"Test Torrent","size":1000000,"progress":0.5,"state":"downloading"}]`))
		} else if r.URL.Path == "/api/v2/auth/login" {
			http.SetCookie(w, &http.Cookie{
				Name:  "SID",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "adminadmin")
	err := client.Login()
	assert.NoError(t, err)

	torrent, err := client.GetTorrentInfo("abc123")
	assert.NoError(t, err)
	assert.NotNil(t, torrent)
	assert.Equal(t, "abc123", torrent.Hash)
}

func TestClient_GetTorrentInfo_NotFound(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/torrents/info" {
			cookie, err := r.Cookie("SID")
			if err != nil || cookie.Value != "test-session-id" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		} else if r.URL.Path == "/api/v2/auth/login" {
			http.SetCookie(w, &http.Cookie{
				Name:  "SID",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok."))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "admin", "adminadmin")
	err := client.Login()
	assert.NoError(t, err)

	torrent, err := client.GetTorrentInfo("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, torrent)
}
