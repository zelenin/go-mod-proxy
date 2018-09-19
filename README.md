# go-mod-proxy

[Go modules proxy](https://tip.golang.org/cmd/go/#hdr-Module_proxy_protocol) boilerplate

## Usage

```go
package main

import (
    "github.com/zelenin/go-mod-proxy"
    "gomod/proxy"
    "log"
    "net/http"
)

func main() {
    // implementation of `gomodproxy.Provider` interface
    provider := proxy.ProviderImplementation()
    
    handler := gomodproxy.New(provider)
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }

    log.Fatal(server.ListenAndServe())
}
```

Then run:

```bash
GOPROXY=http://go-mod-proxy.local:8080 go build ...
```

## Notes

* WIP. Library API can be changed in the future

## Author

[Aleksandr Zelenin](https://github.com/zelenin/), e-mail: [aleksandr@zelenin.me](mailto:aleksandr@zelenin.me)
