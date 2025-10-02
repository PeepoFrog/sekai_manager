package gitres

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type release struct {
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func setUA(req *http.Request) {
	req.Header.Set("User-Agent", "sekaifetch/1.0 (+https://example)")
}

// URLExists does a HEAD (falling back to GET) and treats 200/302/303 as success.
func URLExists(ctx context.Context, client *http.Client, url string) (bool, int, error) {
	{
		req, _ := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
		setUA(req)
		resp, err := client.Do(req)
		if err == nil && resp != nil && resp.StatusCode < 400 {
			resp.Body.Close()
			switch resp.StatusCode {
			case 200, 302, 303:
				return true, resp.StatusCode, nil
			default:
				return false, resp.StatusCode, nil
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	setUA(req)
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200, 302, 303:
		return true, resp.StatusCode, nil
	default:
		return false, resp.StatusCode, nil
	}
}

// FindAssetURL queries /repos/{owner}/{repo}/releases/tags/{tag} (no token)
// and returns the BrowserDownloadURL of the first asset whose name contains nameContains.
func FindAssetURL(ctx context.Context, client *http.Client, owner, repo, tag, nameContains string, verbose bool) (string, error) {
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	setUA(req)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned %d for %s", resp.StatusCode, api)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	for _, a := range rel.Assets {
		if strings.Contains(a.Name, nameContains) {
			if verbose {
				fmt.Println("Selected asset:", a.Name)
			}
			return a.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("no asset matched %q in tag %s", nameContains, tag)
}
