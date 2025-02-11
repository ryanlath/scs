//go:build go1.18
// +build go1.18

package scs

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHijacker(t *testing.T) {
	t.Parallel()

	sessionManager := New()
	sessionManager.Lifetime = 500 * time.Millisecond

	mux := http.NewServeMux()

	mux.HandleFunc("/get", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(http.Hijacker)

		fmt.Fprint(w, ok)
	}))

	ts := newTestServer(t, sessionManager.LoadAndSave(mux))
	defer ts.Close()

	ts.execute(t, "/put")

	_, body := ts.execute(t, "/get")
	if body != "true" {
		t.Errorf("want %q; got %q", "true", body)
	}
}
