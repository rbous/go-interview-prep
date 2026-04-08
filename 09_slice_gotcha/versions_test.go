package slice_gotcha

import (
	"testing"
)

func TestFilterVersions(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-beta", "1.5.0", "3.0.0-rc1", "2.0.0"}

	// Keep only stable versions (no "-" in string)
	stable := FilterVersions(versions, func(v string) bool {
		for _, c := range v {
			if c == '-' {
				return false
			}
		}
		return true
	})

	wantStable := []string{"1.0.0", "1.5.0", "2.0.0"}
	if len(stable) != len(wantStable) {
		t.Fatalf("got %v, want %v", stable, wantStable)
	}
	for i := range wantStable {
		if stable[i] != wantStable[i] {
			t.Errorf("stable[%d] = %q, want %q", i, stable[i], wantStable[i])
		}
	}
}

func TestFilterVersionsOriginalUnchanged(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-beta", "1.5.0", "3.0.0-rc1", "2.0.0"}
	original := make([]string, len(versions))
	copy(original, versions)

	FilterVersions(versions, func(v string) bool {
		return v == "1.0.0"
	})

	for i := range original {
		if versions[i] != original[i] {
			t.Errorf("original modified: versions[%d] = %q, was %q", i, versions[i], original[i])
		}
	}
}

func TestUniqueVersions(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0", "1.0.0", "3.0.0", "2.0.0"}
	unique := UniqueVersions(versions)

	want := []string{"1.0.0", "2.0.0", "3.0.0"}
	if len(unique) != len(want) {
		t.Fatalf("got %v, want %v", unique, want)
	}
	for i := range want {
		if unique[i] != want[i] {
			t.Errorf("unique[%d] = %q, want %q", i, unique[i], want[i])
		}
	}
}

func TestUniqueVersionsOriginalUnchanged(t *testing.T) {
	versions := []string{"a", "b", "a", "c", "b"}
	original := make([]string, len(versions))
	copy(original, versions)

	UniqueVersions(versions)

	for i := range original {
		if versions[i] != original[i] {
			t.Errorf("original modified: versions[%d] = %q, was %q", i, versions[i], original[i])
		}
	}
}
