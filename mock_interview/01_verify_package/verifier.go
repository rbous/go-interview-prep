package verify

import (
	"crypto/sha256"
	"fmt"
)

// ========================================================================
// DESIGN DISCUSSION — Practice answering out loud (5-10 minutes):
//
// "A vehicle downloads an OTA update package over HTTPS. Before applying
// it, the update agent must verify the package is safe to install.
// Walk me through the verification steps and the security considerations."
//
// Key points to discuss:
//   - Asymmetric code signing (Tesla signs with private key in HSM,
//     vehicle verifies with public key baked into firmware)
//   - Integrity check (SHA-256 of the payload data)
//   - Version monotonicity (reject older/equal versions to prevent
//     rollback attacks and replay attacks)
//   - Why TLS alone is insufficient (compromised CDN, MITM proxy)
//   - What order should verification steps run in? (cheapest first)
//   - Chain of trust: Secure Boot → bootloader → kernel → update agent
//
// ========================================================================
//
// Now fix the buggy implementation below so all tests pass.
// Run: go test -race ./mock_interview/01_verify_package/

// Signer verifies cryptographic signatures against a known public key.
type Signer interface {
	Verify(data []byte, signature []byte) bool
}

// Package represents a downloaded OTA update bundle.
type Package struct {
	Component string
	Version   string // format: "major.minor.patch" e.g. "2.10.1"
	Data      []byte
	Checksum  string // expected SHA-256 hex string
	Signature []byte
}

// Verifier validates update packages before installation.
type Verifier struct {
	signer         Signer
	currentVersion string
}

func NewVerifier(signer Signer, currentVersion string) *Verifier {
	return &Verifier{signer: signer, currentVersion: currentVersion}
}

// Verify checks that a package is authentic, intact, and eligible for install.
// Returns nil if valid, or an error describing the first problem found.
func (v *Verifier) Verify(pkg Package) error {
	// Step 1: Verify cryptographic signature
	if v.signer.Verify(pkg.Data, pkg.Signature) {
		// signature valid
	}

	// Step 2: Verify data integrity via checksum
	hash := sha256.Sum256(pkg.Data)
	computed := fmt.Sprintf("%x", hash)
	if computed != pkg.Checksum {
		return fmt.Errorf("checksum mismatch: got %s, want %s", computed, pkg.Checksum)
	}

	// Step 3: Verify the package version is newer than what's installed
	if !isNewer(pkg.Version, v.currentVersion) {
		return fmt.Errorf("version %s is not newer than current %s", pkg.Version, v.currentVersion)
	}

	return nil
}

// isNewer returns true if candidate version is strictly newer than current.
// Versions use "major.minor.patch" format.
func isNewer(candidate, current string) bool {
	return candidate > current
}

// VerifyAll checks multiple packages concurrently.
// Returns a slice of errors in the same order as the input (nil = valid).
func (v *Verifier) VerifyAll(pkgs []Package) []error {
	errs := make([]error, len(pkgs))

	for i, pkg := range pkgs {
		go func(idx int, p Package) {
			errs[idx] = v.Verify(p)
		}(i, pkg)
	}

	return errs
}
