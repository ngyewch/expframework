package assets

import (
	"io/fs"
	"path"
	"sort"
)

func CollectWebAppAssetPaths(webFs fs.FS, assetsDir string) ([]string, error) {
	entries, err := fs.ReadDir(webFs, assetsDir)
	if err != nil {
		return nil, err
	}
	var assetPaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := path.Ext(entry.Name())
		if (ext == ".js") || (ext == ".css") {
			assetPaths = append(assetPaths, "/"+path.Join(assetsDir, entry.Name()))
		}
	}
	sort.Strings(assetPaths)
	return assetPaths, nil
}
