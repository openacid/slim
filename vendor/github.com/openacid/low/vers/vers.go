// Package vers provides version checking functionalities.
package vers

import "strings"
import "github.com/blang/semver"

// VersionGetter defines GetVersion()
//
// Since 0.1.7
type VersionGetter interface {
	GetVersion() string
}

// IsCompatible checks if a verion "ver" satisfies semantic version spec.
//
// `>1.0.0 <2.0.0 || >3.0.0 !4.2.1` would match `1.2.3`, `1.9.9`, `3.1.1`.
// Not `4.2.1`, `2.1.1`
//
// Since 0.1.7
func IsCompatible(ver string, spec []string) bool {

	sp := strings.Join(spec, " || ")

	v, err := semver.Parse(ver)
	if err != nil {
		return false
	}

	chk, err := semver.ParseRange(sp)
	if err != nil {
		return false
	}

	return chk(v)
}
