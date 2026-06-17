package html

import (
	"path"
	"strings"
)

// RewriteAssetURLs rewrites relative href/src attributes in an HTML fragment,
// keeping only the filename and prepending newRoot.
// It ignores absolute URLs, rooted paths, and non-file schemes (data:, mailto:, tel:, javascript:).
func RewriteAssetURLs(htmlStr, newRoot string) string {
	var sb strings.Builder
	lastPos := 0

	for i := 0; i < len(htmlStr); {
		start, end := findNextAttr(htmlStr, i)
		if start == -1 {
			break
		}

		// Write everything up to the match
		sb.WriteString(htmlStr[lastPos:start])

		attrAndEq := htmlStr[start:end]

		// Check what's immediately after the match
		if end >= len(htmlStr) {
			sb.WriteString(attrAndEq)
			lastPos = end
			i = end
			continue
		}

		quote := htmlStr[end : end+1]
		if quote != "\"" && quote != "'" {
			// Unquoted attribute or something else, skip for now to be safe
			sb.WriteString(attrAndEq)
			lastPos = end
			i = end
			continue
		}

		// Find the closing quote
		relIdx := strings.Index(htmlStr[end+1:], quote)
		if relIdx == -1 {
			// No closing quote, skip
			sb.WriteString(attrAndEq)
			lastPos = end
			i = end
			continue
		}

		valStart := end + 1
		valEnd := valStart + relIdx
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

		sb.WriteString(attrAndEq)
		sb.WriteString(quote)
		if shouldRewrite {
			filename := path.Base(oldPath)
			if filename == "." || filename == "/" {
				sb.WriteString(oldPath)
			} else {
				newURL := strings.TrimRight(newRoot, "/") + "/" + filename
				sb.WriteString(newURL)
			}
		} else {
			sb.WriteString(oldPath)
		}
		sb.WriteString(quote)

		lastPos = valEnd + 1
		i = lastPos
	}

	sb.WriteString(htmlStr[lastPos:])
	return sb.String()
}

// findNextAttr finds the next occurrence of src= or href= (case-insensitive)
// with optional spaces around the = and respecting word boundaries at the start.
func findNextAttr(s string, start int) (int, int) {
	for i := start; i < len(s); i++ {
		c := s[i]
		if c != 's' && c != 'S' && c != 'h' && c != 'H' {
			continue
		}

		// \b check (word boundary at the start)
		if i > 0 {
			prev := s[i-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') || prev == '_' {
				continue
			}
		}

		var attrLen int
		if (c == 's' || c == 'S') && i+3 <= len(s) && (s[i:i+3] == "src" || s[i:i+3] == "SRC" || strings.EqualFold(s[i:i+3], "src")) {
			attrLen = 3
		} else if (c == 'h' || c == 'H') && i+4 <= len(s) && (s[i:i+4] == "href" || s[i:i+4] == "HREF" || strings.EqualFold(s[i:i+4], "href")) {
			attrLen = 4
		} else {
			continue
		}

		// Match \s*=\s*
		curr := i + attrLen
		for curr < len(s) && isSpace(s[curr]) {
			curr++
		}
		if curr >= len(s) || s[curr] != '=' {
			continue
		}
		curr++
		for curr < len(s) && isSpace(s[curr]) {
			curr++
		}
		return i, curr
	}
	return -1, -1
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}
