package service

import (
	"errors"
	"testing"
	"time"

	"github.com/iceymoss/inkspace/internal/models"
)

func TestValidateShareLink(t *testing.T) {
	now := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	future := now.Add(time.Hour)
	past := now.Add(-time.Hour)

	tests := []struct {
		name string
		link models.ShareLink
		want error
	}{
		{name: "permanent", link: models.ShareLink{Enabled: true}},
		{name: "not expired", link: models.ShareLink{Enabled: true, ExpiresAt: &future}},
		{name: "expired", link: models.ShareLink{Enabled: true, ExpiresAt: &past}, want: ErrShareExpired},
		{name: "expires exactly now", link: models.ShareLink{Enabled: true, ExpiresAt: &now}, want: ErrShareExpired},
		{name: "disabled", link: models.ShareLink{Enabled: false, ExpiresAt: &future}, want: ErrShareDisabled},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateShareLink(&test.link, now)
			if !errors.Is(err, test.want) {
				t.Fatalf("validateShareLink() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestGenerateShareToken(t *testing.T) {
	first, err := generateShareToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := generateShareToken()
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 24 || first == second {
		t.Fatalf("unexpected tokens %q and %q", first, second)
	}
}
