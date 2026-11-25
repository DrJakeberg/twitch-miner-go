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

	"github.com/Guliveer/twitch-miner-go/internal/constants"
	"github.com/Guliveer/twitch-miner-go/internal/logger"
	"github.com/Guliveer/twitch-miner-go/internal/runtimecfg"
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
	return newAuthenticatorWithCookies(t, rt, map[string]string{"auth-token": "cached-token"})
}

func newAuthenticatorWithCookies(t *testing.T, rt roundTripFunc, cookies map[string]string) *Authenticator {
	t.Helper()
	t.Setenv("TWITCH_AUTH_TOKEN_TESTER", "")
	t.Setenv("TWITCH_PASSWORD_TESTER", "")

	cookieFile := filepath.Join(t.TempDir(), "tester.json")
	jar := NewCookieJar()
	for name, value := range cookies {
		jar.Set(name, value)
	}
	if err := jar.Save(cookieFile); err != nil {
		t.Fatalf("saving cookies: %v", err)
	}

	return &Authenticator{
		username:   "tester",
		cookieJar:  NewCookieJar(),
		cookieFile: cookieFile,
		log:        quietLogger(t),
		httpClient: &http.Client{Transport: rt},
		runtime:    &runtimecfg.Twitch{ClientIDTV: "test-client"},
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

	if !errors.Is(err, errAuthUnavailable) {
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

	if err := a.Login(t.Context()); !errors.Is(err, errAuthUnavailable) {
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

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestLoginKeepsCachedTokenWhenValidateBodyFails(t *testing.T) {
	a := newCachedTokenAuthenticator(t, func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(failingReader{}), Header: make(http.Header)}, nil
	})

	if err := a.Login(t.Context()); !errors.Is(err, errAuthUnavailable) {
		t.Fatalf("expected a body read failure to be transient, got %v", err)
	}
}

func TestLoginKeepsCredentialsWhenRefreshCannotReachTwitch(t *testing.T) {
	a := newAuthenticatorWithCookies(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == constants.TokenURL {
			return nil, &net.DNSError{Err: "i/o timeout", Name: "id.twitch.tv", IsTimeout: true}
		}
		return respond(http.StatusUnauthorized, `{"status":401}`)
	}, map[string]string{"auth-token": "expired-token", "refresh-token": "cached-refresh"})

	if err := a.Login(t.Context()); !errors.Is(err, errAuthUnavailable) {
		t.Fatalf("expected an unreachable refresh to be transient, got %v", err)
	}
}

func TestLoginStoresRefreshedTokensWhenValidationIsUnavailable(t *testing.T) {
	var validations atomic.Int32
	a := newAuthenticatorWithCookies(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == constants.TokenURL {
			return respond(http.StatusOK, `{"access_token":"fresh-access","refresh_token":"fresh-refresh"}`)
		}
		if validations.Add(1) == 1 {
			return respond(http.StatusUnauthorized, `{"status":401}`)
		}
		return respond(http.StatusServiceUnavailable, "")
	}, map[string]string{"auth-token": "expired-token", "refresh-token": "old-refresh"})

	if err := a.Login(t.Context()); !errors.Is(err, errAuthUnavailable) {
		t.Fatalf("expected the post-refresh validation outage to be transient, got %v", err)
	}

	// The old refresh token may already be revoked, so the retry depends on these.
	saved := NewCookieJar()
	if err := saved.Load(a.cookieFile); err != nil {
		t.Fatalf("loading saved cookies: %v", err)
	}
	if got := saved.Get("auth-token"); got != "fresh-access" {
		t.Fatalf("expected the refreshed access token to be saved, got %q", got)
	}
	if got := saved.Get("refresh-token"); got != "fresh-refresh" {
		t.Fatalf("expected the refreshed refresh token to be saved, got %q", got)
	}
}
