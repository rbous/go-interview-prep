package main

import (
	"crypto/sha256"
	"fmt"

	verify "go-interview-prep/mock_interview/01_verify_package"
)

// Simple mock signer for debugging.
type debugSigner struct {
	valid bool
}

func (s *debugSigner) Verify(data []byte, signature []byte) bool {
	return s.valid
}

func hash(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

func main() {
	data := []byte("firmware-v2-binary")

	fmt.Println("=== Package Verifier Debug ===")
	fmt.Println()

	// Case 1: Valid package — should PASS
	fmt.Println("--- Case 1: Valid package ---")
	v := verify.NewVerifier(&debugSigner{valid: true}, "1.0.0")
	err := v.Verify(verify.Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  hash(data),
		Signature: []byte("valid-sig"),
	})
	fmt.Printf("  Result: %v (expected: nil)\n\n", err)

	// Case 2: Invalid signature — should FAIL
	fmt.Println("--- Case 2: Invalid signature ---")
	v2 := verify.NewVerifier(&debugSigner{valid: false}, "1.0.0")
	err = v2.Verify(verify.Package{
		Component: "gateway",
		Version:   "2.0.0",
		Data:      data,
		Checksum:  hash(data),
		Signature: []byte("forged"),
	})
	fmt.Printf("  Result: %v (expected: non-nil error!)\n\n", err)

	// Case 3: Version comparison — "2.10.0" IS newer than "2.9.0"
	fmt.Println("--- Case 3: Multi-digit version ---")
	v3 := verify.NewVerifier(&debugSigner{valid: true}, "2.9.0")
	err = v3.Verify(verify.Package{
		Component: "infotainment",
		Version:   "2.10.0",
		Data:      data,
		Checksum:  hash(data),
		Signature: []byte("sig"),
	})
	fmt.Printf("  Result: %v (expected: nil — 2.10.0 > 2.9.0)\n\n", err)

	// Case 4: VerifyAll — check that results are populated before return
	fmt.Println("--- Case 4: VerifyAll (concurrent) ---")
	v4 := verify.NewVerifier(&debugSigner{valid: true}, "1.0.0")
	pkgs := []verify.Package{
		{Component: "a", Version: "2.0.0", Data: data, Checksum: hash(data), Signature: []byte("sig")},
		{Component: "b", Version: "0.5.0", Data: data, Checksum: hash(data), Signature: []byte("sig")}, // old
		{Component: "c", Version: "2.0.0", Data: data, Checksum: hash(data), Signature: []byte("sig")},
	}
	errs := v4.VerifyAll(pkgs)
	for i, e := range errs {
		fmt.Printf("  pkg[%d] (%s): %v\n", i, pkgs[i].Component, e)
	}
	fmt.Println("  (expected: a=nil, b=version error, c=nil)")
}
