package auth

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Guliveer/twitch-miner-go/internal/logger"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func quietLogger(t *testing.T) *logger.Logger {
	t.Helper()
	cfg := logger.DefaultConfig()
	cfg.Colored = false
	cfg.Level = slog.Level(100)
	l, err := logger.Setup(cfg)
	if err != nil {
		t.Fatalf("logger setup: %v", err)
	}
	return l
}

// newCachedTokenAuthenticator returns an authenticator whose only credential is
// a cached cookie token. skipUnauth turns "fell through every method" into
// ErrSkippedUnauth instead of an interactive device-code login.
func newCachedTokenAuthenticator(t *testing.T, rt roundTripFunc) *Authenticator {
	t.Helper()
	t.Setenv("TWITCH_AUTH_TOKEN_TESTER", "")
	t.Setenv("TWITCH_PASSWORD_TESTER", "")

	cookieFile := filepath.Join(t.TempDir(), "tester.json")
	jar := NewCookieJar()
	jar.Set("auth-token", "cached-token")
	if err := jar.Save(cookieFile); err != nil {
		t.Fatalf("saving cookies: %v", err)
	}

	return &Authenticator{
		username:   "tester",
		cookieJar:  NewCookieJar(),
		cookieFile: cookieFile,
		log:        quietLogger(t),
		httpClient: &http.Client{Transport: rt},
		skipUnauth: true,
	}
}

func respond(status int, body string) (*http.Response, error) {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

func TestLoginKeepsCachedTokenOnNetworkError(t *testing.T) {
	var calls atomic.Int32
	dnsErr := &net.DNSError{Err: "i/o timeout", Name: "id.twitch.tv", IsTimeout: true}
	a := newCachedTokenAuthenticator(t, func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, dnsErr
	})

	err := a.Login(t.Context())

	if !errors.Is(err, errTransientValidation) {
		t.Fatalf("expected a transient validation error, got %v", err)
	}
	var gotDNS *net.DNSError
	if !errors.As(err, &gotDNS) {
		t.Fatalf("expected the DNS cause to be preserved, got %v", err)
	}
	// Retrying belongs to the caller's restart backoff, not to Login itself.
	if n := calls.Load(); n != 1 {
		t.Fatalf("expected a single validation request, got %d", n)
	}
}

func TestLoginKeepsCachedTokenOnServerError(t *testing.T) {
	a := newCachedTokenAuthenticator(t, func(*http.Request) (*http.Response, error) {
		return respond(http.StatusServiceUnavailable, "")
	})

	if err := a.Login(t.Context()); !errors.Is(err, errTransientValidation) {
		t.Fatalf("expected a transient validation error for 503, got %v", err)
	}
}

func TestLoginFallsThroughOnRejectedToken(t *testing.T) {
	a := newCachedTokenAuthenticator(t, func(*http.Request) (*http.Response, error) {
		return respond(http.StatusUnauthorized, `{"status":401,"message":"invalid access token"}`)
	})

	err := a.Login(t.Context())

	if !errors.Is(err, ErrSkippedUnauth) {
		t.Fatalf("expected a rejected token to fall through to the next methods, got %v", err)
	}
}

func TestLoginAcceptsValidCachedToken(t *testing.T) {
	a := newCachedTokenAuthenticator(t, func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, `{"login":"tester","user_id":"42"}`)
	})

	if err := a.Login(t.Context()); err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if got := a.UserID(); got != "42" {
		t.Fatalf("expected user id 42, got %q", got)
	}
}
