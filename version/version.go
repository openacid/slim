package version

const (
	VERSION = "1.0.0"
	/**
	 * never change MAXLEN, it is long enough
	 *
	 * length of version string should be at least
	 * one byte smaller than MAXLEN because having to make room
	 * for delimter in byte array
	 */
	MAXLEN = 16
)
