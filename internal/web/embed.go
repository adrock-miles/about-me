// Package web embeds the server-rendered templates and static assets.
package web

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
)

//go:embed templates/*.html
var Templates embed.FS

//go:embed static
var Static embed.FS

// AssetVersion is a short content hash of every file under static/.
// It changes whenever any embedded asset changes, and stays stable
// otherwise — used as a cache-busting query parameter so browsers and
// CDNs re-fetch styles.css (and friends) after a deploy without us
// having to rename files. Computed once at package init.
var AssetVersion = computeAssetVersion()

func computeAssetVersion() string {
	h := sha256.New()
	err := fs.WalkDir(Static, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(Static, path)
		if err != nil {
			return err
		}
		h.Write([]byte(path))
		h.Write(b)
		return nil
	})
	if err != nil {
		panic("computing asset version: " + err.Error())
	}
	return hex.EncodeToString(h.Sum(nil))[:10]
}
