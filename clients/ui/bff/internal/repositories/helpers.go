package repositories

import (
	"fmt"
	"net/url"
)

func FilterPageValues(values url.Values) url.Values {
	result := url.Values{}

	if v := values.Get("pageSize"); v != "" {
		result.Set("pageSize", v)
	}
	if v := values.Get("orderBy"); v != "" {
		result.Set("orderBy", v)
	}
	if v := values.Get("sortOrder"); v != "" {
		result.Set("sortOrder", v)
	}
	if v := values.Get("nextPageToken"); v != "" {
		result.Set("nextPageToken", v)
	}
	if v := values.Get("name"); v != "" {
		result.Set("name", v)
	}
	if v := values.Get("q"); v != "" {
		result.Set("q", v)
	}
	if v := values.Get("source"); v != "" {
		result.Set("source", v)
	}

	return result
}

func UrlWithParams(url string, values url.Values) string {
	queryString := values.Encode()
	if queryString == "" {
		return url
	}
	return fmt.Sprintf("%s?%s", url, queryString)
}

func UrlWithPageParams(url string, values url.Values) string {
	pageValues := FilterPageValues(values)
	return UrlWithParams(url, pageValues)
}
