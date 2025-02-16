package links

import (
	"fmt"
	"github.com/bheru27/glitzz/config"
	"github.com/bheru27/glitzz/core"
	"github.com/pkg/errors"
	"github.com/thoj/go-ircevent"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const timeout = time.Second * 10

// number of links to read from a single message.
// if this is changed, be sure to change the number of expected links in links_test.go
const linkLimit = 1
const characterLimit = 70

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &links{
		Base:   core.NewBase("links", sender, conf),
		client: &http.Client{Timeout: timeout},
	}
	return rv, nil
}

type links struct {
	core.Base
	client *http.Client
}

func (l *links) HandleEvent(event *irc.Event) {
	if event.Code == "PRIVMSG" {
		links := extractLinks(strings.Fields(event.Message()))
		is_ibip := strings.HasPrefix(event.Message(), "Reporting in! ")
		if is_ibip == false && len(links) >= linkLimit {
			for i, link := range links {
				if i < linkLimit {
					go l.processLink(link, event)
				}
			}
		}
	}
}

func (l *links) processLink(link string, e *irc.Event) {
	title, err := l.getLinkTitle(link)
	if err != nil {
		l.Log.Debug("error getting link title", "link", link, "err", err)
		return
	}
	text := formatResponse(title)
	l.Sender.Reply(e, text)
}

func (l *links) getLinkTitle(link string) (string, error) {
	l.Log.Debug("getting link title", "link", link)
	if strings.HasPrefix(link, "https://twitter.com") {
		link = strings.Replace(link, "twitter.com", "nitter.net", 1)
		l.Log.Debug("Rewrote twitter link", "link", link)
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return "", errors.Wrap(err, "could not create a request")
	}
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	resp, err := l.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http request failed")
	}
	defer resp.Body.Close()
	n, err := html.Parse(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "html parsing failed")
	}
	title, err := findTitle(n)
	if err != nil {
		return "", errors.Wrap(err, "could not find the title")
	}
	title = cleanupTitle(title)

	// Hostname-specific title rendering
	u, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "could not parse url, blame jarboot.")
	}
	switch u.Hostname() {
	case "twitter.com":
		return title, nil
	default:
		if len(title) > characterLimit {
			title = title[:characterLimit-3]
			title += "..."
		}
		return title, nil
	}
}

func cleanupTitle(title string) string {
	re := regexp.MustCompile(`\r?\n`)
	title = re.ReplaceAllString(title, " ")
	title = strings.Join(strings.Fields(title), " ")
	return title
}

func formatResponse(link string) string {
	link = strings.Replace(link, "\n", "", -1)
	return fmt.Sprintf("[ %s ]", link)
}

func isTitleContent(n *html.Node) bool {
	return n.Type == html.TextNode && n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == "title"
}

func findTitle(n *html.Node) (string, error) {
	if isTitleContent(n) {
		return n.Data, nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, err := findTitle(c)
		if err == nil {
			return result, nil
		}
	}
	return "", errors.New("title not found")
}

func extractLinks(arguments []string) []string {
	var links []string
	for _, argument := range arguments {
		if isLink(argument) {
			links = append(links, argument)
		}
	}
	return links
}

func isLink(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")
}
