package tests

import (
	"encoding/base32"
	"strings"
	"testing"
	"time"

	"vessel.dev/vessel/internal/services/oauth"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret1, err := oauth.GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("GenerateTOTPSecret failed: %v", err)
	}
	if len(secret1) == 0 {
		t.Fatal("expected non-empty secret")
	}

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret1)
	if err != nil {
		t.Fatalf("secret is not valid base32: %v", err)
	}
	if len(decoded) != 20 {
		t.Fatalf("expected 20 bytes secret, got %d bytes", len(decoded))
	}

	secret2, err := oauth.GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("GenerateTOTPSecret failed: %v", err)
	}
	if secret1 == secret2 {
		t.Fatal("expected two generated secrets to be distinct")
	}
}

func TestGenerateTOTPQRUri(t *testing.T) {
	uri := oauth.GenerateTOTPQRUri("testuser@example.com", "JBSWY3DPEHPK3PXP")
	expectedPrefix := "otpauth://totp/Vessel:testuser%40example.com"
	if !strings.HasPrefix(uri, expectedPrefix) {
		t.Fatalf("expected prefix %s, got %s", expectedPrefix, uri)
	}
	if !strings.Contains(uri, "secret=JBSWY3DPEHPK3PXP") {
		t.Fatalf("expected secret parameter in uri, got %s", uri)
	}
	if !strings.Contains(uri, "issuer=Vessel") {
		t.Fatalf("expected issuer parameter in uri, got %s", uri)
	}
}

func TestGenerateRecoveryCodes(t *testing.T) {
	codes, err := oauth.GenerateRecoveryCodes(8)
	if err != nil {
		t.Fatalf("GenerateRecoveryCodes failed: %v", err)
	}
	if len(codes) != 8 {
		t.Fatalf("expected 8 codes, got %d", len(codes))
	}

	for _, code := range codes {
		if len(code) != 9 {
			t.Fatalf("expected recovery code length 9, got %d for code %s", len(code), code)
		}
		if code[4] != '-' {
			t.Fatalf("expected hyphen at index 4, got %c in code %s", code[4], code)
		}
	}
}

func TestValidateTOTP(t *testing.T) {
	secret, err := oauth.GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("GenerateTOTPSecret failed: %v", err)
	}

	if oauth.ValidateTOTP(secret, "12345") {
		t.Fatal("expected false for invalid length passcode")
	}
	if oauth.ValidateTOTP("INVALID_BASE32!!!", "123456") {
		t.Fatal("expected false for invalid base32 secret")
	}

	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		t.Fatalf("failed to decode secret: %v", err)
	}

	nowStep := time.Now().Unix() / 30
	validCode := oauth.GenerateTOTPCode(secretBytes, nowStep)
	if !oauth.ValidateTOTP(secret, validCode) {
		t.Fatalf("expected valid code %s to pass validation", validCode)
	}

	futureCode := oauth.GenerateTOTPCode(secretBytes, nowStep+1)
	if !oauth.ValidateTOTP(secret, futureCode) {
		t.Fatalf("expected +1 step code %s to pass validation due to window", futureCode)
	}

	pastCode := oauth.GenerateTOTPCode(secretBytes, nowStep-1)
	if !oauth.ValidateTOTP(secret, pastCode) {
		t.Fatalf("expected -1 step code %s to pass validation due to window", pastCode)
	}

	tooOldCode := oauth.GenerateTOTPCode(secretBytes, nowStep-2)
	if oauth.ValidateTOTP(secret, tooOldCode) {
		t.Fatalf("expected -2 step code %s to fail validation", tooOldCode)
	}
}
