package models

import "testing"

func TestDefaultUserAppearanceResponse(t *testing.T) {
	got := DefaultUserAppearanceResponse()
	if got.UITheme != "magazine" || got.ColorScheme != "system" {
		t.Fatalf("unexpected defaults: %+v", got)
	}
}

func TestUserAppearanceRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		req         UserAppearanceRequest
		wantInvalid bool
	}{
		{name: "system", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "system"}},
		{name: "light", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "light"}},
		{name: "dark", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "dark"}},
		{name: "empty theme", req: UserAppearanceRequest{ColorScheme: "system"}, wantInvalid: true},
		{name: "unavailable theme", req: UserAppearanceRequest{UITheme: "terminal", ColorScheme: "system"}, wantInvalid: true},
		{name: "unknown theme", req: UserAppearanceRequest{UITheme: "unknown", ColorScheme: "system"}, wantInvalid: true},
		{name: "empty color scheme", req: UserAppearanceRequest{UITheme: "magazine"}, wantInvalid: true},
		{name: "unknown color scheme", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "sepia"}, wantInvalid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantInvalid {
				t.Fatalf("Validate() error = %v, wantInvalid = %v", err, tt.wantInvalid)
			}
		})
	}
}
