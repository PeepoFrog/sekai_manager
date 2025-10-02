package downloader

// package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func setUA(req *http.Request) {
	req.Header.Set("User-Agent", "sekaifetch/1.0 (+https://example)")
}

// DownloadToFile streams a URL to a local path. Follows redirects by default via http.Client.
func DownloadToFile(ctx context.Context, client *http.Client, url, path string) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	setUA(req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
