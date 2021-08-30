package urlutil

import "net/url"

func MustParse(rawurl string) *url.URL {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return parsed
}
