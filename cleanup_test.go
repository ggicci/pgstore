package pgstore_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/antonlindstrom/pgstore"
)

func TestCleanup(t *testing.T) {
	db := openDB(t)

	ss, err := pgstore.NewPGStore(db, []byte(secret))
	if err != nil {
		t.Fatal("Failed to get store", err)
	}

	defer ss.Close()
	// Start the cleanup goroutine.
	defer ss.StopCleanup(ss.Cleanup(time.Millisecond * 500))

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("Failed to create request", err)
	}

	session, err := ss.Get(req, "newsess")
	if err != nil {
		t.Fatal("Failed to create session", err)
	}

	// Expire the session.
	session.Options.MaxAge = 1

	m := make(http.Header)
	if err = ss.Save(req, headerOnlyResponseWriter(m), session); err != nil {
		t.Fatal("failed to save session:", err.Error())
	}

	// Give the ticker a moment to run.
	time.Sleep(time.Millisecond * 1500)

	// SELECT expired sessions. We should get a count of zero back.
	var count int
	err = ss.DbPool.QueryRow(context.Background(), "SELECT count(*) FROM http_sessions WHERE expires_on < now()").Scan(&count)
	if err != nil {
		t.Fatalf("failed to select expired sessions from DB: %v", err)
	}

	if count > 0 {
		t.Fatalf("ticker did not delete expired sessions: want 0 got %v", count)
	}
}
