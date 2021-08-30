# go.dpb.io/importshttp

Utilities for hosting Go packages on custom/vanity domains.

 * automatic `go-import` and `go-source` parameters used by Go tooling
 * theming with professional listing and package views by default
 * configurable action links for packages (e.g. `pkg.go.dev`, repository source, custom URLs)
 * used as a server, static site generator, or generic `http.Handler`
 * YAML-based or programmatic configuration

## Usage

For a live example, visit [go.dpb.io](https://go.dpb.io) ([source](https://github.com/dpb587/go.dpb.io); automated via [Cloud Build](https://cloud.google.com/build) and deployed to [Cloud Run](https://cloud.google.com/run)).

To run a server locally, use the [`cmd/http` package](cmd/http)...

```bash
go run go.dpb.io/importshttp/cmd/http \
  -pkg=go.example.com/firstpackage=github.com/golang/go/tree/master \
  -pkg=go.example.com/secondpackage=bitbucket.org/example/secondpackage/src/master
```

To publish a static site with GitHub Pages + Actions, use the [`dpb587/go-importshttp-for-github` template](https://github.com/dpb587/go-importshttp-for-github).

Learn more from the [`examples` directory](examples) and [code documentation](https://pkg.go.dev/go.dpb.io/importshttp).

## License

[MIT License](LICENSE)
