package reqs

import "net/http"

type cookieholder interface {
	setCookieValue(string) error
	setCookieDetails(*http.Cookie)
}

type fragmentholder interface {
	setUrlFragmentValue(string) error
}
