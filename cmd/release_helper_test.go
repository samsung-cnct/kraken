package cmd

import "testing"

func TestCompareReleases(t *testing.T) {
	v1 := "v1.2.3"
	v2 := "v1.3.4"
	v3 := "v0.2"

	// simple test
	c, err := compareReleases(v1, v2)
	if err != nil {
		t.Error("Expected no error in comparison of", v1, "and", v2, "but found error: ", err)
	}

	if c != -1 {
		t.Error("Expected", v1, "to be less than", v2, ", comparison value found to be", c)
	}

	// latest tests
	c, err = compareReleases(latestVersion, v2)
	if err != nil {
		t.Error("Expected no error in comparison of", latestVersion, "and", v2, "but found error: ", err)
	}

	if c != 1 {
		t.Error("Expected", latestVersion, "to be greater than", v2)
	}

	c, err = compareReleases(v1, latestVersion)
	if err != nil {
		t.Error("Expected no error in comparison of", v1, "and", latestVersion, "but found error: ", err)
	}

	if c != -1 {
		t.Error("Expected", v1, "to be less than", latestVersion, ", comparison value found to be", c)
	}

	c, err = compareReleases(latestVersion, latestVersion)
	if err != nil {
		t.Error("Expected no error in comparison of", latestVersion, "and", latestVersion, "but found error: ", err)
	}

	if c != 0 {
		t.Error("Expected", latestVersion, "to be equal to", latestVersion, ", comparison value found to be", c)
	}

	// compare v1 vs v3
	c, err = compareReleases(v1, krakenLibTagToSemver(v3))
	if err != nil {
		t.Error("Expected no error in comparison of", v1, "and", v3, "but found error: ", err)
	}

	if c != 1 {
		t.Error("Expected", v1, "to be greater than", v3, ", comparison value found to be", c)
	}
}

func TestKrakenLibTagToSemver(t *testing.T) {
	v1 := "v0.2"
	v2 := "v1.3.4"

	// test latest first
	v1f := krakenLibTagToSemver(latestVersion)
	if v1f != "latest" {
		t.Error("Expected result to be", latestVersion, "but found", v1f)
	}

	v1f = krakenLibTagToSemver(v1)

	if v1f != "v0.2.0" {
		t.Error("Expected result to be", "v0.2.0", "but found", v1f)
	}

	v1f = krakenLibTagToSemver(v2)

	if v1f != v2 {
		t.Error("Expected result to be", v2, "but found", v1f)
	}

}
