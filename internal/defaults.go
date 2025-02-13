package internal

import (
	"time"

	"github.com/revvim/oxylabs-sdk-go/oxylabs"
)

const (
	DefaultUserAgent oxylabs.UserAgent = oxylabs.UA_DESKTOP
	DefaultDomain    oxylabs.Domain    = oxylabs.DOMAIN_COM

	DefaultStartPage int = 1
	DefaultPages     int = 1

	DefaultLimit_SERP      int = 10
	DefaultLimit_ECOMMERCE int = 48

	SyncBaseUrl  string = "https://realtime.oxylabs.io/v1/queries"
	AsyncBaseUrl string = "https://data.oxylabs.io/v1/queries"
)

var (
	DefaultTimeout      = 50 * time.Second
	DefaultPollInterval = 2 * time.Second
)

// SetDefaultDomain sets the domain parameter if it is not set.
func SetDefaultDomain(domain *oxylabs.Domain) {
	if *domain == "" {
		*domain = oxylabs.DOMAIN_COM
	}
}

// SetDefaultStartPage sets the start_page parameter if it is not set.
func SetDefaultStartPage(startPage *int) {
	if *startPage == 0 {
		*startPage = 1
	}
}

// SetDefaultPages sets the pages parameter if it is not set.
func SetDefaultPages(pages *int) {
	if *pages == 0 {
		*pages = 1
	}
}

// SetDefaultLimit sets the limit parameter if it is not set.
func SetDefaultLimit(limit *int, defaultLimit int) {
	if *limit == 0 {
		*limit = defaultLimit
	}
}

// SetDefaultUserAgent sets the user_agent_type parameter if it is not set.
func SetDefaultUserAgent(userAgent *oxylabs.UserAgent) {
	if *userAgent == "" {
		*userAgent = oxylabs.UA_DESKTOP
	}
}

// SetDefaultHotelOccupancy sets the hotel_occupancy parameter if it is not set.
func SetDefaultHotelOccupancy(ctx oxylabs.ContextOption) {
	if ctx["hotel_occupancy"] == nil {
		ctx["hotel_occupancy"] = 2
	}
}

// SetDefaultSortBy sets the sort_by parameter in the ctx if it is not set.
func SetDefaultSortBy(ctx oxylabs.ContextOption) {
	if ctx["sort_by"] == nil {
		ctx["sort_by"] = "r"
	}
}

// SetDefaultHttpMethod sets the http_method parameter in the ctx if it is not set.
func SetDefaultHttpMethod(ctx oxylabs.ContextOption) {
	if ctx["http_method"] == nil {
		ctx["http_method"] = "get"
	}
}

// SetDefaultContentEncoding sets the content_encoding parameter if it is not set.
func SetDefaultContentEncoding(contentEncoding *string) {
	if *contentEncoding == "" {
		*contentEncoding = "base64"
	}
}
