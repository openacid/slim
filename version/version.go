// Package version defines version of this repo.
//
// It provides backward compatibility checking.
// Most data structure defined in this repo provides with
// serializing/deserializing APIs.
// Serialized data has a version to describe its structure.
//
// Deprecated: will be removed since 1.0.0
package version

const (
	// Latest supported version.
	//
	// On-disk data structure with older version can be read by a newer program.
	VERSION = "1.0.0"

	/**
	 * never change MAXLEN, it is long enough
	 *
	 * length of version string should be at least
	 * one byte smaller than MAXLEN because having to make room
	 * for delimiter in byte array
	 */

	// Fixed size in byte for serialized version.
	//
	// verion is serialized as a string with a trailing '\0'.
	MAXLEN = 16
)

type Version string

type VersionGetter interface {
	GetVersion() Version
}
