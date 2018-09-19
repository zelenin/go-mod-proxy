package gomodproxy

import (
    "errors"
    "io"
    "time"
)

var ErrNotFound = errors.New("package is not found")

type Provider interface {
    Versions(pkgName string) ([]string, error)
    Stat(pkgName string, revision string) (*RevisionInfo, error)
    Latest(pkgName string) (*RevisionInfo, error)
    GoMod(pkgName string, version string) (io.ReadCloser, error)
    Zip(pkgName string, version string) (io.ReadCloser, error)
}

type RevisionInfo struct {
    Version string
    Time    time.Time
}
