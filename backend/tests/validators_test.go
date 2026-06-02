package tests

import (
	"testing"

	"instaforge/internal/validators"
)

func TestValidateInstagramURL_Post(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantType    validators.URLType
		wantErr     bool
	}{
		{"valid post url", "https://instagram.com/p/ABC123def_", validators.URLTypePost, false},
		{"post with www", "https://www.instagram.com/p/ABC123def_", validators.URLTypePost, false},
		{"post no protocol", "instagram.com/p/ABC123def_", validators.URLTypePost, false},
		{"post with query params", "https://instagram.com/p/ABC123def_/?igshid=xyz", validators.URLTypePost, false},
		{"reel url", "https://instagram.com/reel/ABC123def_", validators.URLTypeReel, false},
		{"story url", "https://instagram.com/stories/username/123456789/", validators.URLTypeStory, false},
		{"invalid shortcode (too long)", "https://instagram.com/p/" + string(make([]byte, 100)), validators.URLTypeInvalid, true},
		{"not instagram", "https://youtube.com/watch?v=abc123", validators.URLTypeInvalid, true},
		{"empty url", "", validators.URLTypeInvalid, true},
		{"just path", "/p/ABC123", validators.URLTypeInvalid, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, _, err := validators.ValidateInstagramURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInstagramURL() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if gotType != tt.wantType {
				t.Errorf("ValidateInstagramURL() type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

func TestValidateInstagramUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "test_user", false},
		{"valid with dots", "test.user.123", false},
		{"too short", "a", false},
		{"too long", "averylongusernamethatexceedsthirtycharacters", true},
		{"invalid chars", "user@name", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validators.ValidateInstagramUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInstagramUsername() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "user@example.com", true},
		{"valid with dots", "first.last@example.co.uk", true},
		{"no @", "invalid", false},
		{"no domain", "user@", false},
		{"spaces", "user @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
