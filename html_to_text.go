package main

import (
	"strings"

	"golang.org/x/net/html"
)

// htmlToText converts an HTML string to readable plain text suitable for terminal
// output. It preserves basic block breaks for tags like <p>, <div>, <br>, <li>.
func htmlToText(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return strings.TrimSpace(s)
	}

	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "div", "section", "article", "header", "footer", "tr":
				// paragraph-like block: separate with blank line
				if b.Len() > 0 {
					// avoid adding multiple blank lines
					s := strings.TrimRight(b.String(), " \n")
					b.Reset()
					b.WriteString(s)
					b.WriteString("\n\n")
				}
			case "br":
				b.WriteString("\n")
			case "li":
				// start a new list item with a bullet
				// add newline if needed
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
					b.WriteString("\n")
				}
				b.WriteString("- ")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				if b.Len() > 0 {
					b.WriteString("\n\n")
				}
			}
		}

		if n.Type == html.TextNode {
			// collapse internal whitespace to single spaces
			text := strings.TrimSpace(n.Data)
			if text != "" {
				// if previous char is newline, don't prepend space
				prev := rune(0)
				if b.Len() > 0 {
					prevStr := b.String()
					prev = rune(prevStr[len(prevStr)-1])
				}
				if b.Len() > 0 && prev != '\n' {
					b.WriteString(" ")
				}
				b.WriteString(text)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}

		if n.Type == html.ElementNode {
			// close certain blocks with a blank line
			switch n.Data {
			case "p", "div", "section", "article", "header", "footer", "tr":
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n\n") {
					b.WriteString("\n\n")
				}
			case "li":
				// ensure list item ends with newline
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
					b.WriteString("\n")
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n\n") {
					b.WriteString("\n\n")
				}
			}
		}
	}
	walk(doc)

	out := strings.TrimSpace(b.String())
	// normalize multiple blank lines to max two
	out = strings.ReplaceAll(out, "\n\n\n", "\n\n")
	return out
}

// wrapAndIndent wraps text to the given width and indents each line with the
// provided prefix. The first line also gets the prefix. This is a simple
// word-wrapping implementation suitable for terminal output.
func wrapAndIndent(s string, width int, prefix string) string {
	if width <= 0 {
		return prefix + s
	}
	eff := width - len(prefix)
	if eff <= 20 {
		eff = 20
	}

	var out strings.Builder

	paragraphs := strings.Split(s, "\n\n")
	for pi, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			out.WriteString("\n")
			continue
		}

		// detect list item prefix
		listPrefix := ""
		contIndent := prefix
		if strings.HasPrefix(para, "- ") {
			listPrefix = "- "
			para = strings.TrimPrefix(para, "- ")
			contIndent = prefix + "  "
		}

		words := strings.Fields(para)
		if len(words) == 0 {
			if pi < len(paragraphs)-1 {
				out.WriteString("\n")
			}
			continue
		}

		// first line prefix (includes provided prefix and possible list marker)
		out.WriteString(prefix)
		if listPrefix != "" {
			out.WriteString(listPrefix)
		}

		lineLen := 0
		if listPrefix != "" {
			lineLen = len(listPrefix)
		}
		for i, w := range words {
			if i == 0 {
				out.WriteString(w)
				lineLen += len(w)
				continue
			}
			if lineLen+1+len(w) > eff {
				out.WriteString("\n")
				out.WriteString(contIndent)
				out.WriteString(w)
				lineLen = len(contIndent) - len(prefix) + len(w) // approximate
			} else {
				out.WriteString(" ")
				out.WriteString(w)
				lineLen += 1 + len(w)
			}
		}

		if pi < len(paragraphs)-1 {
			out.WriteString("\n\n")
		}
	}

	return out.String()
}
