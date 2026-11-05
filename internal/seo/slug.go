// Package seo contains the shared public article URL contract.
package seo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// MaxSlugLength is the maximum ASCII length of a persisted article slug.
const MaxSlugLength = 120

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var numericPattern = regexp.MustCompile(`^[0-9]+$`)
var repeatedHyphens = regexp.MustCompile(`-+`)

// IsNumericArticlePath identifies the reserved legacy numeric URL namespace.
func IsNumericArticlePath(value string) bool { return numericPattern.MatchString(value) }

// ValidSlug reports whether value is a canonical, non-numeric single path segment.
func ValidSlug(value string) bool {
	return len(value) <= MaxSlugLength && slugPattern.MatchString(value) && !IsNumericArticlePath(value)
}

// SlugBase converts a title using the pinned local pinyin dictionary. The ID is
// required for titles with no usable characters; numeric titles get a prefix.
func SlugBase(title string, id uint) string {
	var out strings.Builder
	separator := func() {
		if out.Len() > 0 {
			out.WriteByte('-')
		}
	}
	// Do not depend on the library's mutable package-level default options.
	args := pinyin.Args{Style: pinyin.Normal, Fallback: func(rune, pinyin.Args) []string { return nil }}
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			out.WriteRune(r)
		default:
			words := pinyin.SinglePinyin(r, args)
			separator()
			if len(words) > 0 {
				for _, letter := range words[0] {
					if letter >= 'a' && letter <= 'z' {
						out.WriteRune(letter)
					}
				}
				separator()
			}
		}
	}
	value := strings.Trim(out.String(), "-")
	value = repeatedHyphens.ReplaceAllString(value, "-")
	if value == "" {
		value = fmt.Sprintf("article-%d", id)
	}
	if IsNumericArticlePath(value) {
		value = "article-" + value
	}
	return strings.TrimRight(value[:min(len(value), MaxSlugLength)], "-")
}

// SlugCandidate returns base, base-2, base-3, ... without exceeding the limit.
func SlugCandidate(base string, attempt uint) string {
	if attempt < 2 {
		return base
	}
	suffix := fmt.Sprintf("-%d", attempt)
	return strings.TrimRight(base[:min(len(base), MaxSlugLength-len(suffix))], "-") + suffix
}

// ArticlePath returns the canonical relative path for a validated persisted slug.
func ArticlePath(slug string) string { return "/blog/" + slug }
