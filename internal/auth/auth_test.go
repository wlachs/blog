package auth_test

import (
	"github.com/wlchs/blog/internal/auth"
	"testing"
)

// TestHashStringValid tests whether hashing method works correctly.
func TestHashStringValid(t *testing.T) {
	t.Parallel()

	h1, _ := auth.HashString("test")

	if len(h1) != 60 {
		t.Errorf("hash length mismatch")
	}

	h2, _ := auth.HashString("test")
	if h1 == h2 {
		t.Errorf("hashes should use random salt")
	}
}

// TestHashStringInvalidLong tests how the hashing method handles too long input.
func TestHashStringInvalidLong(t *testing.T) {
	t.Parallel()

	in := "$2y$10$2/rIv3UPAU0llQpJeM2aiuiL8BNl3OlTs/uVSIGiSm6QwF2q2ddo21234567890123"
	_, err := auth.HashString(in)

	if err == nil {
		t.Errorf("should not be able to hash too long strings")
	}
}

// TestCompareStringWithHash compares a generated hash with its plaintext counterpart.
func TestCompareStringWithHash(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		plaintext string
		hash      string
		match     bool
	}{
		"#1: Should match":     {plaintext: "Test", hash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.", match: true},
		"#2: Should match":     {plaintext: "Test1", hash: "$2y$10$o7ZqUckxyaZAS31yFfBhTutbo3cWUQkvsdnVikvhrn69.c5kG0/TS", match: true},
		"#3: Should not match": {plaintext: "Test", hash: "$2y$10$tnhWpPEURh779WSzhA/G9eCU1edd/Y29V9X9IoW8qmUSdOmZXLJIG", match: false},
		"#4: Should not match": {plaintext: "Test1", hash: "$2y$10$2/rIv3UPAU0llQpJeM2aiuiL8BNl3OlTs/uVSIGiSm6QwF2q2ddo2", match: false},
	}

	for scenario, tc := range tt {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			match := auth.CompareStringWithHash(tc.plaintext, tc.hash)

			if match != tc.match {
				t.Errorf("failed hash comparison: %s - %s", tc.plaintext, tc.hash)
			}
		})
	}
}
