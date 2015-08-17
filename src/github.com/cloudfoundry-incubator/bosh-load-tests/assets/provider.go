package assets

import (
	"path/filepath"
)

type provider struct {
	baseDir string
}

type Provider interface {
	FullPath(path string) string
}

func NewProvider(baseDir string) Provider {
	return &provider{
		baseDir: baseDir,
	}
}

func (p *provider) FullPath(path string) string {
	return filepath.Join(p.baseDir, path)
}
