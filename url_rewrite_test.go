package html_test

import (
	"testing"

	. "github.com/tinywasm/html"
)

func TestRewriteAssetURLs(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		newRoot string
		want    string
	}{
		{
			name:    "simple src double quotes",
			html:    `<img src="images/logo.png">`,
			newRoot: "/static",
			want:    `<img src="/static/logo.png">`,
		},
		{
			name:    "simple href single quotes",
			html:    `<link rel='stylesheet' href='css/style.css'>`,
			newRoot: "/assets",
			want:    `<link rel='stylesheet' href='/assets/style.css'>`,
		},
		{
			name:    "ignore absolute URL http",
			html:    `<script src="http://example.com/lib.js"></script>`,
			newRoot: "/static",
			want:    `<script src="http://example.com/lib.js"></script>`,
		},
		{
			name:    "ignore rooted path",
			html:    `<img src="/already/rooted.png">`,
			newRoot: "/static",
			want:    `<img src="/already/rooted.png">`,
		},
		{
			name:    "ignore data URL",
			html:    `<img src="data:image/png;base64,abc">`,
			newRoot: "/static",
			want:    `<img src="data:image/png;base64,abc">`,
		},
		{
			name:    "ignore empty src",
			html:    `<img src="">`,
			newRoot: "/static",
			want:    `<img src="">`,
		},
		{
			name:    "handle nested quotes if mismatched",
			html:    `<img src="it's-a-picture.png">`,
			newRoot: "/static",
			want:    `<img src="/static/it's-a-picture.png">`,
		},
		{
			name:    "mismatched quotes should be ignored",
			html:    `<img src="wrong'>`,
			newRoot: "/static",
			want:    `<img src="wrong'>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RewriteAssetURLs(tt.html, tt.newRoot)
			if got != tt.want {
				t.Errorf("RewriteAssetURLs(%s) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
