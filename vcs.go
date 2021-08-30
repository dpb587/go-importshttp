package importshttp

// VCS is any go mod-supported version control system.
type VCS string

const (
	UnknownVCS    VCS = ""
	BazaarVCS     VCS = "bzr"
	FossilVCS     VCS = "fossil"
	GitVCS        VCS = "git"
	MercurialVCS  VCS = "hg"
	ModVCS        VCS = "mod"
	SubversionVCS VCS = "svn"
)
