package mcardhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ReposClient contains a list of repos related github api methods.
type ReposClient interface {
	List(user string) ([]*Repos, error)
	ListLanguages(url string) (map[string]int, error)
}

// Repos represents a subset of GitHub repository response object.
type Repos struct {
	ID          *int64  `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
	Description *string `json:"description,omitempty"`

	// Language URL
	LanguagesURL *string `json:"languages_url,omitempty"`
}

// addRequestHeaders is helper method to append custom headers to the Request Object.
func (c *V1Client) addRequestHeaders(r *http.Request) *http.Request {
	r.Header.Set("Content-type", "application/json")
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("User-Agent", "Mcard Client")

	// Setting this field prevents re-use of TCP connections between requests to the same hosts, as if
	// Transport.DisableKeepAlives were set.
	r.Close = true

	return r
}

// List returns list of repos associated to user
func (c *V1Client) List(user string) ([]*Repos, error) {
	if user == "" {
		return nil, errors.New("user not provided")
	}
	url := c.baseURL
	url.Path = fmt.Sprintf("users/%v/repos", user)

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	res, err := c.httpu.MakeRequest(c.addRequestHeaders(req))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var repos []*Repos

	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// ListLanguages provides the list of languages associated with the repo
func (c *V1Client) ListLanguages(url string) (map[string]int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	res, err := c.httpu.MakeRequest(c.addRequestHeaders(req))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	languages := make(map[string]int)

	if err := json.NewDecoder(res.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}
