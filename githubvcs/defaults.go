package githubvcs

// DefaultServer is the default provider for GitHub repositories.
var DefaultServer = "https://github.com"

// DefaultRef is the ref/branch used when an explicit one is not provided. To avoid confusion, ref should always be
// configured and this default not relied upon.
//
// Previously "master" prior to https://github.blog/changelog/2020-10-01-the-default-branch-for-newly-created-repositories-is-now-main/.
var DefaultRef = "main"
