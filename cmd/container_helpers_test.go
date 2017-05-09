package cmd




import(
	"testing"
	"os"
)




func TestAddEnvironmentVarIfNotEmpty(t *testing.T) {
	// consider random generator in the future for testing.
	keyOne := "TEST_KEY_1"
	keyTwo := "TEST_KEY_2"
	keyThree := "TEST_KEY_3"

	valOne := "Value1"
	valTwo := "Value2"
	valThree := ""

	os.Setenv(keyOne, valOne)
	os.Setenv(keyTwo, valTwo)
	os.Setenv(keyThree, valThree)

	envs := []string{"ANSIBLE_NOCOLOR=True", "DISPLAY_SKIPPED_HOSTS=0"}

	envs = appendIfValueNotEmpty(envs, keyOne)
	envs = appendIfValueNotEmpty(envs, keyTwo)
	envs = appendIfValueNotEmpty(envs, keyThree)

	if len(envs) != 4 {
		t.Error("For", envs, "expetected length to be", 4, "got", len(envs))
	}

	badEnvVar := keyThree+"="
	for _, s := range envs {
		if s == badEnvVar {
			t.Error("For", envs, "did not expect to find key", badEnvVar, "but found it.")
		}
	}

	goodEnvVar := keyTwo+"="+valTwo
	goodEnvRes := false
	for _, s := range envs {
		if s == goodEnvVar {
			goodEnvRes = true
		}
	}

	if !goodEnvRes {
		t.Error("For", envs, "expected to find key", goodEnvVar, "but not found.")
	}
}
