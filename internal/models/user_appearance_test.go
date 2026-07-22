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
		{name: "magazine system", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "system"}},
		{name: "magazine light", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "light"}},
		{name: "magazine dark", req: UserAppearanceRequest{UITheme: "magazine", ColorScheme: "dark"}},
		{name: "terminal system", req: UserAppearanceRequest{UITheme: "terminal", ColorScheme: "system"}},
		{name: "terminal light", req: UserAppearanceRequest{UITheme: "terminal", ColorScheme: "light"}},
		{name: "terminal dark", req: UserAppearanceRequest{UITheme: "terminal", ColorScheme: "dark"}},
		{name: "empty theme", req: UserAppearanceRequest{ColorScheme: "system"}, wantInvalid: true},
		{name: "cozy unavailable", req: UserAppearanceRequest{UITheme: "cozy", ColorScheme: "system"}, wantInvalid: true},
		{name: "swiss unavailable", req: UserAppearanceRequest{UITheme: "swiss", ColorScheme: "system"}, wantInvalid: true},
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
