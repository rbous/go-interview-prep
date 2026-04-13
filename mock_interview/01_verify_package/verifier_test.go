package verify

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

// --- Test helpers (do NOT modify) ---

type mockSigner struct {
	valid bool
}

func (m *mockSigner) Verify(data []byte, signature []byte) bool {
	return m.valid
}

func checksum(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

// --- Tests ---

func TestVerifyValidPackage(t *testing.T) {
	data := []byte("firmware-v2-binary")
	v := NewVerifier(&mockSigner{valid: true}, "1.0.0")

	err := v.Verify(Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  checksum(data),
		Signature: []byte("valid-sig"),
	})
	if err != nil {
		t.Errorf("valid package should pass verification: %v", err)
	}
}

func TestRejectsInvalidSignature(t *testing.T) {
	data := []byte("firmware-v2-binary")
	v := NewVerifier(&mockSigner{valid: false}, "1.0.0")

	err := v.Verify(Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  checksum(data),
		Signature: []byte("forged-sig"),
	})
	if err == nil {
		t.Error("expected error for invalid signature, got nil — a forged package would be installed!")
	}
}

func TestRejectsChecksumMismatch(t *testing.T) {
	v := NewVerifier(&mockSigner{valid: true}, "1.0.0")

	err := v.Verify(Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      []byte("firmware"),
		Checksum:  "0000000000000000000000000000000000000000000000000000000000000000",
		Signature: []byte("sig"),
	})
	if err == nil {
		t.Error("expected checksum mismatch error")
	}
}

func TestVersionComparisonMultiDigit(t *testing.T) {
	// 2.10.0 is newer than 2.9.0, but string comparison says otherwise
	// because "1" < "9" at the character level.
	data := []byte("fw")
	v := NewVerifier(&mockSigner{valid: true}, "2.9.0")

	err := v.Verify(Package{
		Component: "infotainment",
		Version:   "2.10.0",
		Data:      data,
		Checksum:  checksum(data),
		Signature: []byte("sig"),
	})
	if err != nil {
		t.Errorf("2.10.0 should be accepted as newer than 2.9.0: %v", err)
	}
}

func TestRejectsOlderVersion(t *testing.T) {
	data := []byte("fw")
	v := NewVerifier(&mockSigner{valid: true}, "3.0.0")

	err := v.Verify(Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  checksum(data),
		Signature: []byte("sig"),
	})
	if err == nil {
		t.Error("expected error — older version should be rejected")
	}
}

func TestRejectsEqualVersion(t *testing.T) {
	// Installing the same version again could be a replay attack.
	data := []byte("fw")
	v := NewVerifier(&mockSigner{valid: true}, "2.0.0")

	err := v.Verify(Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  checksum(data),
		Signature: []byte("sig"),
	})
	if err == nil {
		t.Error("expected error — equal version should be rejected (replay protection)")
	}
}

func TestVerifyAllCompletes(t *testing.T) {
	data := []byte("fw")
	v := NewVerifier(&mockSigner{valid: true}, "1.0.0")

	pkgs := []Package{
		{Component: "a", Version: "2.0.0", Data: data, Checksum: checksum(data), Signature: []byte("sig")},
		{Component: "b", Version: "0.5.0", Data: data, Checksum: checksum(data), Signature: []byte("sig")}, // old version
		{Component: "c", Version: "2.0.0", Data: data, Checksum: checksum(data), Signature: []byte("sig")},
	}

	errs := v.VerifyAll(pkgs)

	if errs[0] != nil {
		t.Errorf("package 'a' should be valid: %v", errs[0])
	}
	if errs[1] == nil {
		t.Error("package 'b' should fail — version 0.5.0 is older than current 1.0.0")
	}
	if errs[2] != nil {
		t.Errorf("package 'c' should be valid: %v", errs[2])
	}
}
