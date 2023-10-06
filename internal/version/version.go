package version

import "strings"

const (
	versionMajor = "2"
	versionMinor = "0"
	versionPatch = "0"
)

func VersionString() string {
	return "v" + strings.Join([]string{versionMajor, versionMinor, versionPatch}, ".")
}
