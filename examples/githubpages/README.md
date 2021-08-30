# go-importshttp-for-github

This repository manages the custom remote imports for Go packages through [GitHub](https://github.com/). It relies on:

 * [`config.yaml` file](config.yaml) for site and package settings;
 * [`dpb587/go-importshttp` package](https://github.com/dpb587/go-importshttp) for generating a Go-compatible, static site for the packages;
 * [GitHub Actions](https://github.com/features/actions) for rebuilding the site on `gh-pages` whenever the settings change; and
 * [GitHub Pages](https://pages.github.com/) for serving the packages site.

Useful links:

 * [`dpb587/go-importshttp-for-github`](https://github.com/dpb587/go-importshttp-for-github) (original repository template)
 * [Configuring a custom domain for your GitHub Pages site](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site)
