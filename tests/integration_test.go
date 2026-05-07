package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestApplicationEndToEnd(t *testing.T) {
	client, baseURL := startTestApp(t)

	externalID := "3f93df6d-ff51-4740-9d27-fc6b2f30281c"
	body := fmt.Sprintf(`{
		"external_id": %q,
		"name": "Jane Doe",
		"email": "jane@example.com",
		"date_of_birth": "1990-01-02T03:04:05Z"
	}`, externalID)
	expected := personResponse{
		ExternalID:  externalID,
		Name:        "Jane Doe",
		Email:       "jane@example.com",
		DateOfBirth: "1990-01-02T03:04:05Z",
	}

	createResp := createPerson(t, client, baseURL, body)
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected POST /save status 201, got %d", createResp.StatusCode)
	}
	if location := createResp.Header.Get("Location"); location != "/"+externalID {
		t.Fatalf("expected Location %q, got %q", "/"+externalID, location)
	}
	if created := decodePersonResponse(t, createResp); created != expected {
		t.Fatalf("expected created response %+v, got %+v", expected, created)
	}

	getResp := getPerson(t, client, baseURL, externalID)
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected GET /{id} status 200, got %d", getResp.StatusCode)
	}
	if saved := decodePersonResponse(t, getResp); saved != expected {
		t.Fatalf("expected saved response %+v, got %+v", expected, saved)
	}

	duplicateResp := createPerson(t, client, baseURL, body)
	defer duplicateResp.Body.Close()
	if duplicateResp.StatusCode != http.StatusConflict {
		t.Fatalf("expected duplicate POST /save status 409, got %d", duplicateResp.StatusCode)
	}

	unknownResp := getPerson(t, client, baseURL, "f965f31d-ef51-446d-b5b0-41c7ab77e1d0")
	defer unknownResp.Body.Close()
	if unknownResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected GET unknown UUID status 404, got %d", unknownResp.StatusCode)
	}

	invalidResp := getPerson(t, client, baseURL, "not-a-uuid")
	defer invalidResp.Body.Close()
	if invalidResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected GET invalid UUID status 400, got %d", invalidResp.StatusCode)
	}
}

type personResponse struct {
	ExternalID  string `json:"external_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
}

func startTestApp(t *testing.T) (*http.Client, string) {
	t.Helper()

	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}

	port := freeTCPPort(t)
	dbPath := filepath.Join(t.TempDir(), "app.db")
	baseURL := "http://127.0.0.1:" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/api")
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"HTTP_PORT="+port,
		"DB_PATH="+dbPath,
	)
	// go run starts a compiled child binary; killing the process group cleans up both.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start app: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		cancel()
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Log("timed out waiting for app process to exit")
		}
		if t.Failed() {
			t.Logf("stdout:\n%s", stdout.String())
			t.Logf("stderr:\n%s", stderr.String())
		}
	})

	client := &http.Client{Timeout: 2 * time.Second}
	waitForHealth(t, client, baseURL)

	return client, baseURL
}

func freeTCPPort(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on free port: %v", err)
	}
	defer listener.Close()

	return fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
}

func waitForHealth(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("app did not become healthy before timeout")
}

func createPerson(t *testing.T, client *http.Client, baseURL string, body string) *http.Response {
	t.Helper()

	return doRequest(t, client, http.MethodPost, baseURL+"/save", body)
}

func getPerson(t *testing.T, client *http.Client, baseURL string, externalID string) *http.Response {
	t.Helper()

	return doRequest(t, client, http.MethodGet, baseURL+"/"+externalID, "")
}

func decodePersonResponse(t *testing.T, resp *http.Response) personResponse {
	t.Helper()

	var decoded personResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("decode person response: %v", err)
	}

	return decoded
}

func doRequest(t *testing.T, client *http.Client, method string, url string, body string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("send %s %s: %v", method, url, err)
	}

	return resp
}
