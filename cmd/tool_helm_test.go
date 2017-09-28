package cmd

import (
	"testing"
)

func TestRemovePatchVersion(t *testing.T) {
	var res string
	var err error

	// error cases
	errorCases := []string{
		"",
		"v",
		"v1",
		"1.0.4",
		"v1",
		"v2.0.",
		"v4.6.3.2",
	}

	for _, ec := range errorCases {
		if res, err = removePatchVersion(ec); err == nil {
			t.Error("For error case", ec, "was not expecting to successfully remove patch version, result was", res)
		}
	}

	// success cases
	successCases := []string{
		"v1.6.7",
		"v121.5434.1253",
	}

	expectedOutput := []string{
		"v1.6",
		"v121.5434",
	}

	for i, sc := range successCases {
		res, err = removePatchVersion(sc)
		if err != nil {
			t.Error("For success case", sc, "was not expecting to error:", err)
		}
		if res != expectedOutput[i] {
			t.Error("For success case", sc, "was expecting", expectedOutput[i], "but found", res)
		}
	}
}
