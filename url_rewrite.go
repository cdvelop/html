//go:build !wasm

package html

import (
	"path"
	"regexp"
	"strings"
)

// attrStartRegex matches (src|href) followed by an optional space and = and optional space.
// It uses \b to ensure we match whole attributes only.
var attrStartRegex = regexp.MustCompile(`(?i)\b(src|href)\s*=\s*`)

// RewriteAssetURLs rewrites relative href/src attributes in an HTML fragment,
// keeping only the filename and prepending newRoot.
// It ignores absolute URLs, rooted paths, and non-file schemes (data:, mailto:, tel:, javascript:).
// SSR/backend-only: this is excluded from wasm builds (its only caller, assetmin, is backend).
func RewriteAssetURLs(htmlStr, newRoot string) string {
	indices := attrStartRegex.FindAllStringIndex(htmlStr, -1)
	if indices == nil {
		return htmlStr
	}

	var sb strings.Builder
	lastPos := 0

	for _, idx := range indices {
		// Write everything up to the match
		sb.WriteString(htmlStr[lastPos:idx[0]])

		// The attribute being processed (src or href)
		attrAndEq := htmlStr[idx[0]:idx[1]]

		// Check what's immediately after the =
		if idx[1] >= len(htmlStr) {
			sb.WriteString(attrAndEq)
			lastPos = idx[1]
			continue
		}

		quote := htmlStr[idx[1] : idx[1]+1]
		if quote != "\"" && quote != "'" {
			// Unquoted attribute or something else, skip for now to be safe
			sb.WriteString(attrAndEq)
			lastPos = idx[1]
			continue
		}

		// Find the closing quote
		endIdx := strings.Index(htmlStr[idx[1]+1:], quote)
		if endIdx == -1 {
			// No closing quote, skip
			sb.WriteString(attrAndEq)
			lastPos = idx[1]
			continue
		}

		valStart := idx[1] + 1
		valEnd := valStart + endIdx
		oldPath := htmlStr[valStart:valEnd]

		shouldRewrite := true
		if oldPath == "" {
			shouldRewrite = false
		} else if strings.HasPrefix(oldPath, "http://") ||
			strings.HasPrefix(oldPath, "https://") ||
			strings.HasPrefix(oldPath, "//") ||
			strings.HasPrefix(oldPath, "/") {
			shouldRewrite = false
		} else {
			lowerPath := strings.ToLower(oldPath)
			if strings.HasPrefix(lowerPath, "data:") ||
				strings.HasPrefix(lowerPath, "mailto:") ||
				strings.HasPrefix(lowerPath, "tel:") ||
				strings.HasPrefix(lowerPath, "javascript:") {
				shouldRewrite = false
			}
		}

		if shouldRewrite {
			filename := path.Base(oldPath)
			if filename == "." || filename == "/" {
				sb.WriteString(attrAndEq)
				sb.WriteString(quote)
				sb.WriteString(oldPath)
				sb.WriteString(quote)
			} else {
				newURL := strings.TrimRight(newRoot, "/") + "/" + filename
				sb.WriteString(attrAndEq)
				sb.WriteString(quote)
				sb.WriteString(newURL)
				sb.WriteString(quote)
			}
		} else {
			sb.WriteString(attrAndEq)
			sb.WriteString(quote)
			sb.WriteString(oldPath)
			sb.WriteString(quote)
		}

		lastPos = valEnd + 1
	}

	sb.WriteString(htmlStr[lastPos:])
	return sb.String()
}
