package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blakesmith/ar"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

// ExtractFirstMatch opens a .deb (ar archive), finds data.tar.{gz,xz,zst},
// and extracts the first file whose basename matches any name in basenames.
// The extracted file is written to destDir with its basename; returns full path.
func ExtractFirstMatch(debPath string, basenames []string, destDir string) (string, error) {
	f, err := os.Open(debPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	arR := ar.NewReader(f)
	var dataBuf bytes.Buffer
	var dataName string

	for {
		hdr, err := arR.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(hdr.Name, "data.tar") {
			dataName = hdr.Name
			if _, err := io.Copy(&dataBuf, arR); err != nil {
				return "", err
			}
			break
		}
	}
	if dataBuf.Len() == 0 {
		return "", fmt.Errorf("data.tar.* not found in %s", debPath)
	}

	tr, closer, err := openTarStream(dataName, bytes.NewReader(dataBuf.Bytes()))
	if err != nil {
		return "", err
	}
	defer func() {
		if closer != nil {
			_ = closer.Close()
		}
	}()

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

nextEntry:
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if !h.FileInfo().Mode().IsRegular() {
			continue
		}
		base := filepath.Base(h.Name)
		for _, want := range basenames {
			if base == want {
				outPath := filepath.Join(destDir, want)
				tmp := outPath + ".partial"
				if err := writeFileFromTar(tmp, tr, h.FileInfo().Mode()); err != nil {
					return "", err
				}
				if err := os.Rename(tmp, outPath); err != nil {
					return "", err
				}
				return outPath, nil
			}
		}
		// continue scanning
		continue nextEntry
	}

	return "", fmt.Errorf("none of %v found in data.tar", basenames)
}

func writeFileFromTar(path string, r io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = io.Copy(f, r); err != nil {
		return err
	}
	// ensure executable
	_ = os.Chmod(path, mode|0o111)
	return nil
}

// openTarStream detects compression from filename suffix (gz, xz, zst) and returns a tar.Reader.
// The returned io.Closer (if non-nil) should be closed by the caller.
func openTarStream(name string, src io.Reader) (*tar.Reader, io.Closer, error) {
	switch {
	case strings.HasSuffix(name, ".gz"):
		gzr, err := gzip.NewReader(src)
		if err != nil {
			return nil, nil, fmt.Errorf("gzip open: %w", err)
		}
		return tar.NewReader(gzr), gzr, nil
	case strings.HasSuffix(name, ".xz"):
		xzr, err := xz.NewReader(src)
		if err != nil {
			return nil, nil, fmt.Errorf("xz open: %w", err)
		}
		// xz.Reader doesn't implement io.Closer
		return tar.NewReader(xzr), nil, nil
	case strings.HasSuffix(name, ".zst"), strings.HasSuffix(name, ".zstd"):
		zr, err := zstd.NewReader(src)
		if err != nil {
			return nil, nil, fmt.Errorf("zstd open: %w", err)
		}
		return tar.NewReader(zr), zr.IOReadCloser(), nil
	default:
		// Try gzip by magic header; fallback to plain tar
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, src); err != nil {
			return nil, nil, err
		}
		b := buf.Bytes()
		if len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b {
			gzr, err := gzip.NewReader(bytes.NewReader(b))
			if err != nil {
				return nil, nil, fmt.Errorf("gzip(magic) open: %w", err)
			}
			return tar.NewReader(gzr), gzr, nil
		}
		return tar.NewReader(bytes.NewReader(b)), nil, nil
	}
}
