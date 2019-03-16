package hashutils

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func HashHTMLBody(b []byte) string {
	f := sha512.New()

	doc, _ := html.Parse(strings.NewReader(string(b)))
	bn, err := getBody(doc)
	if err != nil {
		return "error"
	}
	body, err := renderNode(bn)
	if err != nil {
		return "error"
	}
	dom := html.NewTokenizer(strings.NewReader(body))
	previous := dom.Token()
loopDom:
	for {
		tt := dom.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDom
		case tt == html.StartTagToken:
			previous = dom.Token()
		case tt == html.TextToken:
			if previous.Data == "script" {
				continue
			}
			txt := strings.TrimSpace(html.UnescapeString(string(dom.Text())))
			if len(txt) > 0 {
				_, _ = io.WriteString(f, txt)
			}
		}
	}
	s := base64.URLEncoding.EncodeToString(f.Sum(nil))
	return s
}

func getBody(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("missing <body>")
}

func renderNode(n *html.Node) (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, n)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
