package fetch

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func getClientWithCookie(userID string, cookies ...*http.Cookie) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies = append(cookies, &http.Cookie{
		Name:  "token",
		Value: userID,
	})

	jar.SetCookies(&url.URL{
		Scheme: "http",
		Host:   "localhost:" + os.Getenv("PORT"),
	}, cookies)

	c := &http.Client{
		Jar: jar,
	}

	return c, nil
}
