package gomodproxy

import (
    "encoding/json"
    "github.com/zelenin/go-router"
    "io"
    "net/http"
)

const (
    moduleKey   = "module"
    revisionKey = "revision"
    versionKey  = "version"

    modulePattern  = `.+`
    encodedPattern = `[\w\d+\-\.~]+`
)

func New(provider Provider) http.Handler {
    rtr := router.New()

    rtr.Get(`/{`+moduleKey+`:`+modulePattern+`}/@v/list`, VersionsHandler{provider})
    rtr.Get(`/{`+moduleKey+`:`+modulePattern+`}/@v/{`+revisionKey+`:`+encodedPattern+`}.info`, StatHandler{provider})
    rtr.Get(`/{`+moduleKey+`:`+modulePattern+`}/@v/latest`, LatestHandler{provider})
    rtr.Get(`/{`+moduleKey+`:`+modulePattern+`}/@v/{`+versionKey+`:`+encodedPattern+`}.mod`, GoModHandler{provider})
    rtr.Get(`/{`+moduleKey+`:`+modulePattern+`}/@v/{`+versionKey+`:`+encodedPattern+`}.zip`, ZipHandler{provider})

    rtr.Pipe("/", checkModuleName)

    return rtr
}

func checkModuleName(next http.Handler) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        params := router.Params(req)

        moduleName := params.Get(moduleKey)

        if moduleName == "" {
            res.WriteHeader(404)

            return
        }

        next.ServeHTTP(res, req)
    })
}

// Versions lists all known versions with the given prefix.
// Pseudo-versions are not included.
type VersionsHandler struct {
    provider Provider
}

func (handler VersionsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    params := router.Params(req)

    versions, err := handler.provider.Versions(params.Get(moduleKey))
    if err != nil {
        if err == ErrNotFound {
            res.WriteHeader(404)
        } else {
            res.WriteHeader(500)
        }

        // @todo log

        return
    }

    res.Header().Set("Content-type", "text/plain")

    for _, version := range versions {
        res.Write([]byte(version))
    }
}

// Stat returns information about the revision rev.
// A revision can be any identifier known to the underlying service: commit hash, branch, tag, and so on.
type StatHandler struct {
    provider Provider
}

func (handler StatHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    params := router.Params(req)

    revisionInfo, err := handler.provider.Stat(params.Get(moduleKey), params.Get(revisionKey))
    if err != nil {
        if err == ErrNotFound {
            res.WriteHeader(404)
        } else {
            res.WriteHeader(500)
        }

        // @todo log

        return
    }

    res.Header().Set("Content-type", "application/json")

    err = json.NewEncoder(res).Encode(revisionInfo)
    if err != nil {
        res.WriteHeader(500)
        // @todo log

        return
    }
}

// Latest returns the latest revision on the default branch,
// whatever that means in the underlying source code repository.
// It is only used when there are no tagged versions.
type LatestHandler struct {
    provider Provider
}

func (handler LatestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    params := router.Params(req)

    revisionInfo, err := handler.provider.Latest(params.Get(moduleKey))
    if err != nil {
        if err == ErrNotFound {
            res.WriteHeader(404)
        } else {
            res.WriteHeader(500)
        }

        // @todo log

        return
    }

    res.Header().Set("Content-type", "application/json")

    err = json.NewEncoder(res).Encode(revisionInfo)
    if err != nil {
        res.WriteHeader(500)
        // @todo log

        return
    }
}

// GoMod returns the go.mod file for the given version.
type GoModHandler struct {
    provider Provider
}

func (handler GoModHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    params := router.Params(req)

    goMod, err := handler.provider.GoMod(params.Get(moduleKey), params.Get(versionKey))
    if err != nil {
        if err == ErrNotFound {
            res.WriteHeader(404)
        } else {
            res.WriteHeader(500)
        }

        // @todo log

        return
    }

    defer goMod.Close()

    res.Header().Set("Content-type", "text/plain")

    io.Copy(res, goMod)
}

// Zip downloads a zip file for the given version to a new file in a given temporary directory.
// It returns the name of the new file.
// The caller should remove the file when finished with it.
type ZipHandler struct {
    provider Provider
}

func (handler ZipHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    params := router.Params(req)

    zip, err := handler.provider.Zip(params.Get(moduleKey), params.Get(versionKey))
    if err != nil {
        if err == ErrNotFound {
            res.WriteHeader(404)
        } else {
            res.WriteHeader(500)
        }

        // @todo log

        return
    }

    defer zip.Close()

    res.Header().Set("Content-type", "application/zip")

    io.Copy(res, zip)
}
