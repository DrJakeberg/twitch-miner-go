package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Guliveer/twitch-miner-go/internal/version"
)

const (
	releaseURL = "https://api.github.com/repos/Guliveer/twitch-miner-go/releases/latest"
	repoURL    = "https://github.com/Guliveer/twitch-miner-go"
	timeout    = 5 * time.Second
)

// UpdateInfo holds the result of an update check.
type UpdateInfo struct {
	Available bool
	Latest    string
	URL       string
	IsGitRepo bool
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// CheckForUpdate checks GitHub for a newer release.
func CheckForUpdate(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	return checkWithURL(ctx, currentVersion, releaseURL)
}

func checkWithURL(ctx context.Context, currentVersion, url string) (*UpdateInfo, error) {
	if currentVersion == "dev" {
		return &UpdateInfo{Available: false}, nil
	}

	current, err := version.Parse(currentVersion)
	if err != nil {
		return &UpdateInfo{Available: false}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	latestStr := strings.TrimPrefix(release.TagName, "v")
	latest, err := version.Parse(latestStr)
	if err != nil {
		return nil, fmt.Errorf("parse remote version %q: %w", release.TagName, err)
	}

	info := &UpdateInfo{
		Available: version.Compare(latest, current) > 0,
		Latest:    latestStr,
		URL:       repoURL,
		IsGitRepo: isGitRepo(),
	}
	return info, nil
}

// FormatNotification returns the user-facing update message.
func FormatNotification(info *UpdateInfo, currentVersion string) string {
	if !info.Available {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n══════════════════════════════════════════════════════════\n")
	fmt.Fprintf(&b, "  🔔 New version available: v%s (current: v%s)\n", info.Latest, currentVersion)
	b.WriteString("\n")
	if info.IsGitRepo {
		b.WriteString("  Update:\n")
		b.WriteString("    git pull && ./run.sh\n")
		b.WriteString("\n")
		b.WriteString("  Or download manually:\n")
	} else {
		b.WriteString("  Download the latest version:\n")
	}
	fmt.Fprintf(&b, "    %s\n", info.URL)
	b.WriteString("══════════════════════════════════════════════════════════\n")
	return b.String()
}

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return err == nil
}
