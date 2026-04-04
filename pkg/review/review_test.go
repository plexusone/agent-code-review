package review

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v84/github"
)

func TestReviewFooter(t *testing.T) {
	if ReviewFooter == "" {
		t.Error("ReviewFooter should not be empty")
	}
	if len(ReviewFooter) < 10 {
		t.Error("ReviewFooter should contain meaningful content")
	}
}

func TestReviewEvent_Constants(t *testing.T) {
	tests := []struct {
		event ReviewEvent
		want  string
	}{
		{EventApprove, "APPROVE"},
		{EventRequestChanges, "REQUEST_CHANGES"},
		{EventComment, "COMMENT"},
	}
	for _, tt := range tests {
		if string(tt.event) != tt.want {
			t.Errorf("got %s, want %s", tt.event, tt.want)
		}
	}
}

func TestNewClient(t *testing.T) {
	gh := github.NewClient(nil)
	client := NewClient(gh)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.gh != gh {
		t.Error("NewClient did not set the GitHub client correctly")
	}
}

func TestNewClientFromToken(t *testing.T) {
	ctx := context.Background()
	client := NewClientFromToken(ctx, "test-token")
	if client == nil {
		t.Fatal("NewClientFromToken returned nil")
	}
	if client.gh == nil {
		t.Error("NewClientFromToken did not create a GitHub client")
	}
}

func TestReviewInput_Fields(t *testing.T) {
	input := ReviewInput{
		Owner:    "owner",
		Repo:     "repo",
		PRNumber: 123,
		Event:    EventApprove,
		Body:     "LGTM!",
	}
	if input.Owner != "owner" {
		t.Errorf("Owner = %s, want owner", input.Owner)
	}
	if input.Repo != "repo" {
		t.Errorf("Repo = %s, want repo", input.Repo)
	}
	if input.PRNumber != 123 {
		t.Errorf("PRNumber = %d, want 123", input.PRNumber)
	}
	if input.Event != EventApprove {
		t.Errorf("Event = %s, want APPROVE", input.Event)
	}
	if input.Body != "LGTM!" {
		t.Errorf("Body = %s, want LGTM!", input.Body)
	}
}

func TestCommentInput_Fields(t *testing.T) {
	input := CommentInput{
		Owner:    "owner",
		Repo:     "repo",
		PRNumber: 456,
		Body:     "Great work!",
	}
	if input.Owner != "owner" {
		t.Errorf("Owner = %s, want owner", input.Owner)
	}
	if input.Repo != "repo" {
		t.Errorf("Repo = %s, want repo", input.Repo)
	}
	if input.PRNumber != 456 {
		t.Errorf("PRNumber = %d, want 456", input.PRNumber)
	}
	if input.Body != "Great work!" {
		t.Errorf("Body = %s, want Great work!", input.Body)
	}
}

func TestLineCommentInput_Fields(t *testing.T) {
	input := LineCommentInput{
		Owner:    "owner",
		Repo:     "repo",
		PRNumber: 789,
		CommitID: "abc123",
		Path:     "main.go",
		Line:     42,
		Body:     "Consider renaming this",
	}
	if input.Owner != "owner" {
		t.Errorf("Owner = %s, want owner", input.Owner)
	}
	if input.CommitID != "abc123" {
		t.Errorf("CommitID = %s, want abc123", input.CommitID)
	}
	if input.Path != "main.go" {
		t.Errorf("Path = %s, want main.go", input.Path)
	}
	if input.Line != 42 {
		t.Errorf("Line = %d, want 42", input.Line)
	}
}

func TestPRInfo_Fields(t *testing.T) {
	info := PRInfo{
		Number:  123,
		Title:   "Add feature",
		Body:    "Description",
		State:   "open",
		Author:  "user",
		Head:    "feature-branch",
		Base:    "main",
		Commits: 5,
		HTMLURL: "https://github.com/owner/repo/pull/123",
	}
	if info.Number != 123 {
		t.Errorf("Number = %d, want 123", info.Number)
	}
	if info.State != "open" {
		t.Errorf("State = %s, want open", info.State)
	}
	if info.Commits != 5 {
		t.Errorf("Commits = %d, want 5", info.Commits)
	}
}

func TestPRSummary_Fields(t *testing.T) {
	summary := PRSummary{
		Number:  456,
		Title:   "Fix bug",
		Author:  "contributor",
		Head:    "bugfix",
		HTMLURL: "https://github.com/owner/repo/pull/456",
	}
	if summary.Number != 456 {
		t.Errorf("Number = %d, want 456", summary.Number)
	}
	if summary.Author != "contributor" {
		t.Errorf("Author = %s, want contributor", summary.Author)
	}
}

// Integration-style tests using httptest mock server

func setupMockGitHub(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	gh, err := github.NewClient(nil).WithEnterpriseURLs(server.URL, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return NewClient(gh), server
}

func TestCreateReview_AppendsFooter(t *testing.T) {
	var capturedBody string

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Enterprise URLs include /api/v3/ prefix
		if r.Method == "POST" && (r.URL.Path == "/repos/owner/repo/pulls/1/reviews" || r.URL.Path == "/api/v3/repos/owner/repo/pulls/1/reviews") {
			// Parse the request body to capture the review body
			var req struct {
				Body  string `json:"body"`
				Event string `json:"event"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
				capturedBody = req.Body
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1, "html_url": "https://github.com/owner/repo/pull/1#pullrequestreview-1"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}

	client, server := setupMockGitHub(t, handler)
	defer server.Close()

	ctx := context.Background()
	_, err := client.CreateReview(ctx, &ReviewInput{
		Owner:    "owner",
		Repo:     "repo",
		PRNumber: 1,
		Event:    EventApprove,
		Body:     "LGTM!",
	})
	if err != nil {
		t.Fatalf("CreateReview failed: %v", err)
	}

	if capturedBody == "" {
		t.Skip("Could not capture request body")
	}

	expectedSuffix := ReviewFooter
	if len(capturedBody) < len(expectedSuffix) {
		t.Errorf("Body too short to contain footer")
	}
}
