package main

import (
  "fmt"
  "log"
  "strings"
  "net/http"
  "golang.org/x/net/html"
)

// Fetch returns the body of URL and a slice of URLs found on that page.
type Fetcher interface {
  Fetch(url string) (body string, urls []string, err error)
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
  // Iterate over all of the Token's attributes until we find an "href"
  for _, a := range t.Attr {
    if a.Key == "href" {
      href = a.Val
      ok = true
    }
  }

  return
}

// Crawl uses fetcher to recursively crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
  if depth <= 0 {
    return
  }
  
  body, urls, err := fetcher.Fetch(url)
  
  if err != nil {
    fmt.Println(err)
    return
  }
  
	fmt.Printf("found: %s %q\n", url, body)

	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}

	return  
}

type URLFetcher struct {
  
}

func (f URLFetcher) Fetch(url string) (string, []string, error) {
  resp, err := http.Get(url)
  
  if err != nil {
    return "", nil, err
  }
  
  body := resp.Body
  defer body.Close()

  // collect the urls
  var urls []string
  z := html.NewTokenizer(body)
  for {
    tt := z.Next()

    switch {
    case tt == html.ErrorToken:
      // End of the document, we're done
      return "", urls, nil

    case tt == html.StartTagToken:
      t := z.Token()

      // Check if the token is an <a> tag
      isAnchor := t.Data == "a"
      if !isAnchor {
        continue
      }

      // Extract the href value, if there is one
      ok, url := getHref(t)
      if !ok {
        continue
      }

      // Make sure the url begins with http**
      hasProto := strings.Index(url, "http") == 0
      if hasProto {
        urls = append(urls, url)
      }
    }
  }
}

func main() {
	Crawl("http://golang.org/", 4, URLFetcher{})
}