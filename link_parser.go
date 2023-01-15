package link_parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// Parse all html links in input data
// HTML link must be with `href` property
// returns array of Link structs and error in case
// of parsing incoming data error
func Parse(rdr io.Reader) ([]Link, error) {
	doc, err := html.Parse(rdr)
	if err != nil {
		return nil, err
	}
	var res []Link
	q := []*html.Node{doc}
	for len(q) > 0 {
		crnt := q[0]
		if crnt.FirstChild != nil {
			fc := crnt.FirstChild
			q = append(q, fc)

			for fc.NextSibling != nil {
				q = append(q, fc.NextSibling)
				fc = fc.NextSibling
			}
		}
		if crnt.Type == html.ElementNode && crnt.Data == "a" {
			res = append(res, *parseLink(crnt))
		}
		q = q[1:]
	}

	return res, nil
}

// Parse html.Node struct with type html.ElementNode to Link struct
func parseLink(htmlN *html.Node) *Link {
	var href string
	for _, attr := range htmlN.Attr {
		if attr.Key == "href" {
			href = attr.Val
		}
	}
	lnk := Link{
		Href: href,
		Text: strContent(*htmlN),
	}
	return &lnk
}

// Returns string from all children of the node
func strContent(nd html.Node) string {
	var stk []*html.Node
	var res strings.Builder
	if nd.Type == html.ElementNode && nd.FirstChild != nil {
		stk = append(stk, nd.FirstChild)
	}

	for len(stk) > 0 {
		crnt := stk[len(stk)-1]
		stk = stk[1:]
		if crnt.Type == html.TextNode {
			res.WriteString(crnt.Data)
		}
		if crnt.Type == html.ElementNode {
			res.WriteString(strContent(*crnt))
		}
		if crnt != nil && crnt.NextSibling != nil {
			stk = append(stk, crnt.NextSibling)
		}
	}

	return strings.Join(strings.Fields(res.String()), " ")
}
