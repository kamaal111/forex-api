package integration

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	waitTimeout     = 10 * time.Second
	pollInterval    = 100 * time.Millisecond
	shutdownTimeout = 5 * time.Second
	firestoreHost   = "127.0.0.1:8080"
)

type ServerProcess struct {
	cmd     *exec.Cmd
	Port    int
	baseURL string
}

type TestContext struct {
	Server    *ServerProcess
	DB        *firestore.Client
	Ctx       context.Context
	ProjectID string
}

func NewTestContext() *TestContext {
	projectID := fmt.Sprintf("forex-api-test-%d", time.Now().UnixNano())
	return &TestContext{
		Ctx:       context.Background(),
		ProjectID: projectID,
	}
}

func (tc *TestContext) Setup(port int) error {
	var err error
	tc.DB, err = firestore.NewClient(tc.Ctx, tc.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %w", err)
	}

	tc.Server, err = StartServer(tc.ProjectID, port)
	if err != nil {
		tc.DB.Close()
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (tc *TestContext) Teardown() {
	if tc.Server != nil {
		tc.Server.Stop()
	}
	if tc.DB != nil {
		tc.DB.Close()
	}
}

func (tc *TestContext) ClearCollection(collectionName string) error {
	collection := tc.DB.Collection(collectionName)
	docs, err := collection.Documents(tc.Ctx).GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		return nil
	}

	batch := tc.DB.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	_, err = batch.Commit(tc.Ctx)
	return err
}

func StartServer(projectID string, port int) (*ServerProcess, error) {
	if port == 0 {
		var err error
		port, err = getFreePort()
		if err != nil {
			return nil, fmt.Errorf("failed to get free port: %w", err)
		}
	}

	serverAddress := fmt.Sprintf("127.0.0.1:%d", port)
	baseURL := fmt.Sprintf("http://%s", serverAddress)

	projectRoot := getProjectRoot()

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GCP_PROJECT_ID=%s", projectID),
		fmt.Sprintf("SERVER_ADDRESS=%s", serverAddress),
		fmt.Sprintf("FIRESTORE_EMULATOR_HOST=%s", firestoreHost),
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	server := &ServerProcess{
		cmd:     cmd,
		Port:    port,
		baseURL: baseURL,
	}

	if err := waitForServer(baseURL, waitTimeout); err != nil {
		server.Stop()
		return nil, fmt.Errorf("server failed to become ready: %w", err)
	}

	return server, nil
}

func (s *ServerProcess) Stop() {
	if s.cmd == nil || s.cmd.Process == nil {
		return
	}

	s.cmd.Process.Signal(os.Interrupt)
	done := make(chan error, 1)
	go func() {
		done <- s.cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(shutdownTimeout):
		s.cmd.Process.Kill()
	}
}

func (s *ServerProcess) BaseURL() string {
	return s.baseURL
}

func (s *ServerProcess) GetLatest(base, symbols string) (*http.Response, error) {
	url := fmt.Sprintf("%s/v1/rates/latest", s.baseURL)
	if base != "" || symbols != "" {
		url += "?"
		if base != "" {
			url += fmt.Sprintf("base=%s", base)
		}
		if symbols != "" {
			if base != "" {
				url += "&"
			}
			url += fmt.Sprintf("symbols=%s", symbols)
		}
	}

	return http.Get(url)
}

func waitForServer(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: pollInterval}

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/v1/rates/latest")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(pollInterval)
	}

	return fmt.Errorf("server did not respond within %v", timeout)
}

func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return wd
}
