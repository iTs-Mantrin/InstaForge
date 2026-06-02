package validators

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Instagram URL patterns
var (
	postURLPattern   = regexp.MustCompile(`^(https?://)?(www\.)?(instagram\.com)/(p|tv)/([A-Za-z0-9_-]+)`)
	reelURLPattern   = regexp.MustCompile(`^(https?://)?(www\.)?(instagram\.com)/(reel)/([A-Za-z0-9_-]+)`)
	storyURLPattern  = regexp.MustCompile(`^(https?://)?(www\.)?(instagram\.com)/(stories)/([A-Za-z0-9_.-]+)/([0-9]+)`)
	usernamePattern  = regexp.MustCompile(`^[A-Za-z0-9_.]{1,30}$`)
	shortcodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,30}$`)
)

// URLType represents the type of Instagram URL detected.
type URLType string

const (
	URLTypePost   URLType = "post"
	URLTypeReel   URLType = "reel"
	URLTypeStory  URLType = "story"
	URLTypeInvalid URLType = "invalid"
)

// ValidateInstagramURL validates and returns the type of an Instagram URL.
func ValidateInstagramURL(rawURL string) (URLType, string, error) {
	// Clean and normalize URL
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return URLTypeInvalid, "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Normalize host
	host := strings.TrimPrefix(parsedURL.Hostname(), "www.")
	if host != "instagram.com" {
		return URLTypeInvalid, "", fmt.Errorf("not an Instagram URL")
	}

	fullPath := rawURL

	if matches := postURLPattern.FindStringSubmatch(fullPath); len(matches) > 5 {
		shortcode := matches[5]
		if !shortcodePattern.MatchString(shortcode) {
			return URLTypePost, "", fmt.Errorf("invalid post shortcode")
		}
		return URLTypePost, shortcode, nil
	}

	if matches := reelURLPattern.FindStringSubmatch(fullPath); len(matches) > 5 {
		shortcode := matches[5]
		if !shortcodePattern.MatchString(shortcode) {
			return URLTypeReel, "", fmt.Errorf("invalid reel shortcode")
		}
		return URLTypeReel, shortcode, nil
	}

	if matches := storyURLPattern.FindStringSubmatch(fullPath); len(matches) > 6 {
		username := matches[5]
		storyID := matches[6]
		if !usernamePattern.MatchString(username) {
			return URLTypeStory, "", fmt.Errorf("invalid username in story URL")
		}
		return URLTypeStory, fmt.Sprintf("%s_%s", username, storyID), nil
	}

	return URLTypeInvalid, "", fmt.Errorf("unsupported Instagram URL format")
}

// ValidateInstagramUsername validates an Instagram username.
func ValidateInstagramUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 1 || len(username) > 30 {
		return fmt.Errorf("username must be between 1 and 30 characters")
	}
	if !usernamePattern.MatchString(username) {
		return fmt.Errorf("username contains invalid characters (alphanumeric, underscore, period only)")
	}
	return nil
}

// ExtractUsernameFromURL extracts a username from an Instagram profile URL.
func ExtractUsernameFromURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", fmt.Errorf("no username in URL")
	}

	username := parts[0]
	if !usernamePattern.MatchString(username) {
		return "", fmt.Errorf("invalid username format")
	}

	return username, nil
}

// IsValidEmail performs basic email validation.
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
