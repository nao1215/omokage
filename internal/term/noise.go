package term

import (
	"regexp"
	"strings"
)

var (
	// frontmatter matches a leading YAML front-matter block (Hugo/Jekyll posts),
	// whose paths and metadata are not prose and would otherwise flood the
	// candidates with "images", "jpg", "post", and the like.
	frontmatter = regexp.MustCompile(`(?s)\A\s*---\n.*?\n---\n`)
	// markdownLink keeps a link/image's visible text and drops its URL, so a link
	// target never contributes URL fragments (https, com, jp) as candidates.
	markdownLink = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]*\)`)
	htmlTag      = regexp.MustCompile(`<[^>]*>`)
	url          = regexp.MustCompile(`https?://\S+`)
	htmlEntity   = regexp.MustCompile(`&[a-zA-Z]+;`)
)

// stripNoise removes non-prose that survives code stripping but should never
// become a term candidate or feed an alias bridge: YAML front matter, link/image
// URLs, raw URLs, HTML tags, and HTML entities. It is deterministic and adds no
// dependency. Visible link text is preserved so real terms inside links still
// count.
func stripNoise(text string) string {
	text = frontmatter.ReplaceAllString(text, "")
	text = markdownLink.ReplaceAllString(text, "$1")
	text = url.ReplaceAllString(text, " ")
	text = htmlTag.ReplaceAllString(text, " ")
	text = htmlEntity.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
