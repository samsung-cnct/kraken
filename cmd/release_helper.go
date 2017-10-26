package cmd

import (
	"strings"

	"fmt"

	"github.com/coreos/go-semver/semver"
)

// compareReleases expects semver versions are being compared
// compareReleases tests if versionA is less than, equal to, or greater than versionB,
// returning -1, 0, or +1 respectively.
// -2 is default of error cases
func compareReleases(versionA string, versionB string) (int, error) {
	versionA = strings.Replace(versionA, "v", "", -1)
	versionB = strings.Replace(versionB, "v", "", -1)

	if versionA == latestVersion && versionB == latestVersion {
		return 0, nil
	}

	if versionA == latestVersion && versionB != latestVersion {
		return 1, nil
	}

	if versionA != latestVersion && versionB == latestVersion {
		return -1, nil
	}

	v1, err := semver.NewVersion(versionA)
	if err != nil {
		return -2, err
	}

	v2, err := semver.NewVersion(versionB)
	if err != nil {
		return -2, err
	}

	return v1.Compare(*v2), nil
}

func krakenLibTagToSemver(tag string) string {
	if tag == latestVersion {
		return tag
	}

	tagSlice := strings.Split(tag, ".")

	switch len(tagSlice) {
	case 3:
		return tag
	case 2:
		return fmt.Sprintf("%s.0", tag)
	default:
		return tag
	}
}
