package ecommerce

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/revvim/oxylabs-sdk-go/internal"
	"github.com/revvim/oxylabs-sdk-go/oxylabs"
)

// Accepted parameters for context options in google shopping.
var AcceptedSortByParameters = []string{
	"r",
	"p",
	"rv",
	"pd",
}

// GoogleShoppingUrlOpts contains all the query parameters available for google shopping.
type GoogleShoppingUrlOpts struct {
	UserAgent         oxylabs.UserAgent
	Render            oxylabs.Render
	CallbackUrl       string
	GeoLocation       string
	Parse             bool
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
}

// checkParameterValidity checks validity of ScrapeGoogleShoppingUrl parameters.
func (opt *GoogleShoppingUrlOpts) checkParameterValidity() error {
	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Render != "" && !oxylabs.IsRenderValid(opt.Render) {
		return fmt.Errorf("invalid render parameter: %v", opt.Render)
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// ScrapeGoogleShoppingUrl scrapes google shopping via Oxylabs E-Commerce API with google_shopping as source.
func (c *EcommerceClient) ScrapeGoogleShoppingUrl(
	url string,
	opts ...*GoogleShoppingUrlOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeGoogleShoppingUrlCtx(ctx, url, opts...)
}

// ScrapeGoogleShoppingUrlCtx scrapes google shopping via Oxylabs E-Commerce API with google_shopping as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *EcommerceClient) ScrapeGoogleShoppingUrlCtx(
	ctx context.Context,
	url string,
	opts ...*GoogleShoppingUrlOpts,
) (*Resp, error) {
	// Check validity of url.
	err := internal.ValidateUrl(url, "shopping.google")
	if err != nil {
		return nil, err
	}

	// Prepare options.
	opt := &GoogleShoppingUrlOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Set defaults.
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err = opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare payload.
	payload := map[string]interface{}{
		"source":          oxylabs.GoogleShoppingUrl,
		"url":             url,
		"user_agent_type": opt.UserAgent,
		"render":          opt.Render,
		"callback_url":    opt.CallbackUrl,
		"geo_location":    opt.GeoLocation,
		"parse":           opt.Parse,
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, opt.Parse, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GoogleShoppingSearchOpts contains all the query parameters available for google shopping search.
type GoogleShoppingSearchOpts struct {
	Domain            oxylabs.Domain
	StartPage         int
	Pages             int
	Locale            oxylabs.Locale
	ResultsLanguage   string
	GeoLocation       string
	UserAgent         oxylabs.UserAgent
	Render            oxylabs.Render
	CallbackURL       string
	Parse             bool
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
	Context           []func(oxylabs.ContextOption)
}

// checkParameterValidity checks validity of ScrapeGoogleShoppingSearch parameters.
func (opt *GoogleShoppingSearchOpts) checkParameterValidity(ctx oxylabs.ContextOption) error {
	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Render != "" && !oxylabs.IsRenderValid(opt.Render) {
		return fmt.Errorf("invalid render parameter: %v", opt.Render)
	}

	if opt.Pages <= 0 || opt.StartPage <= 0 {
		return fmt.Errorf("pages and start_page parameters must be greater than 0")
	}

	if ctx["sort_by"] != nil && !internal.InList(ctx["sort_by"].(string), AcceptedSortByParameters) {
		return fmt.Errorf("invalid sort_by parameter: %v", ctx["sort_by"])
	}

	if (ctx["min_price"] != nil || ctx["max_price"] != nil) &&
		(ctx["min_price"].(int) < 0 || ctx["max_price"].(int) < 0) {
		return fmt.Errorf("min and max prices should be greater than 0")
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// ScrapeGoogleShoppingSearch scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_search as source.
func (c *EcommerceClient) ScrapeGoogleShoppingSearch(
	query string,
	opts ...*GoogleShoppingSearchOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeGoogleShoppingSearchCtx(ctx, query, opts...)
}

// ScrapeGoogleShoppingSearchCtx scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_search as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *EcommerceClient) ScrapeGoogleShoppingSearchCtx(
	ctx context.Context,
	query string,
	opts ...*GoogleShoppingSearchOpts,
) (*Resp, error) {
	// Prepare options.
	opt := &GoogleShoppingSearchOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Initialize the context map and apply each provided context modifier function.
	context := make(oxylabs.ContextOption)
	for _, modifier := range opt.Context {
		modifier(context)
	}

	// Set defaults.
	internal.SetDefaultSortBy(context)
	internal.SetDefaultPages(&opt.Pages)
	internal.SetDefaultDomain(&opt.Domain)
	internal.SetDefaultStartPage(&opt.StartPage)
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err := opt.checkParameterValidity(context)
	if err != nil {
		return nil, err
	}

	// Prepare payload with common parameters.
	payload := map[string]interface{}{
		"source":           oxylabs.GoogleShoppingSearch,
		"domain":           opt.Domain,
		"query":            query,
		"start_page":       opt.StartPage,
		"pages":            opt.Pages,
		"locale":           opt.Locale,
		"results_language": opt.ResultsLanguage,
		"geo_location":     opt.GeoLocation,
		"user_agent_type":  opt.UserAgent,
		"render":           opt.Render,
		"callback_url":     opt.CallbackURL,
		"parse":            opt.Parse,
		"context": []map[string]interface{}{
			{
				"key":   "nfpr",
				"value": context["nfpr"],
			},
			{
				"key":   "sort_by",
				"value": context["sort_by"],
			},
			{
				"key":   "min_price",
				"value": context["min_price"],
			},
			{
				"key":   "max_price",
				"value": context["max_price"],
			},
		},
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, opt.Parse, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GoogleShoppingProductOpts contains all the query parameters available for google shopping product.
type GoogleShoppingProductOpts struct {
	Domain            oxylabs.Domain
	Locale            oxylabs.Locale
	ResultsLanguage   string
	GeoLocation       string
	UserAgent         oxylabs.UserAgent
	Render            oxylabs.Render
	CallbackURL       string
	Parse             bool
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
}

// checkParameterValidity checks validity of ScrapeGoogleShoppingProduct parameters.
func (opt *GoogleShoppingProductOpts) checkParameterValidity() error {
	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Render != "" && !oxylabs.IsRenderValid(opt.Render) {
		return fmt.Errorf("invalid render parameter: %v", opt.Render)
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// ScrapeGoogleShoppingProduct scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_product as source.
func (c *EcommerceClient) ScrapeGoogleShoppingProduct(
	query string,
	opts ...*GoogleShoppingProductOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeGoogleShoppingProductCtx(ctx, query, opts...)
}

// ScrapeGoogleShoppingProductCtx scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_product as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *EcommerceClient) ScrapeGoogleShoppingProductCtx(
	ctx context.Context,
	query string,
	opts ...*GoogleShoppingProductOpts,
) (*Resp, error) {
	// Prepare options.
	opt := &GoogleShoppingProductOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Set defaults.
	internal.SetDefaultDomain(&opt.Domain)
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err := opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare payload with common parameters.
	payload := map[string]interface{}{
		"source":           oxylabs.GoogleShoppingProduct,
		"domain":           opt.Domain,
		"query":            query,
		"locale":           opt.Locale,
		"results_language": opt.ResultsLanguage,
		"geo_location":     opt.GeoLocation,
		"user_agent_type":  opt.UserAgent,
		"render":           opt.Render,
		"callback_url":     opt.CallbackURL,
		"parse":            opt.Parse,
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, opt.Parse, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GoogleShoppingPricingOpts contains all the query parameters available for google shopping pricing.
type GoogleShoppingPricingOpts struct {
	Domain            oxylabs.Domain
	StartPage         int
	Pages             int
	Locale            oxylabs.Locale
	ResultsLanguage   string
	GeoLocation       string
	UserAgent         oxylabs.UserAgent
	Render            oxylabs.Render
	CallbackURL       string
	Parse             bool
	ParseInstructions *map[string]interface{}
	PollInterval      time.Duration
}

// checkParameterValidity checks validity of ScrapeGoogleShoppingPricing parameters.
func (opt *GoogleShoppingPricingOpts) checkParameterValidity() error {
	if !oxylabs.IsUserAgentValid(opt.UserAgent) {
		return fmt.Errorf("invalid user agent parameter: %v", opt.UserAgent)
	}

	if opt.Render != "" && !oxylabs.IsRenderValid(opt.Render) {
		return fmt.Errorf("invalid render parameter: %v", opt.Render)
	}

	if opt.Pages <= 0 || opt.StartPage <= 0 {
		return fmt.Errorf("pages and start_page parameters must be greater than 0")
	}

	if opt.ParseInstructions != nil {
		if err := oxylabs.ValidateParseInstructions(opt.ParseInstructions); err != nil {
			return fmt.Errorf("invalid parse instructions: %w", err)
		}
	}

	return nil
}

// ScrapeGoogleShoppingPricing scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_pricing as source.
func (c *EcommerceClient) ScrapeGoogleShoppingPricing(
	query string,
	opts ...*GoogleShoppingPricingOpts,
) (*Resp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.DefaultTimeout)
	defer cancel()

	return c.ScrapeGoogleShoppingPricingCtx(ctx, query, opts...)
}

// ScrapeGoogleShoppingPricingCtx scrapes google shopping via Oxylabs E-Commerce API
// with google_shopping_pricing as source.
// The provided context allows customization of the HTTP req, including setting timeouts.
func (c *EcommerceClient) ScrapeGoogleShoppingPricingCtx(
	ctx context.Context,
	query string,
	opts ...*GoogleShoppingPricingOpts,
) (*Resp, error) {
	// Prepare options.
	opt := &GoogleShoppingPricingOpts{}
	if len(opts) > 0 && opts[len(opts)-1] != nil {
		opt = opts[len(opts)-1]
	}

	// Set defaults.
	internal.SetDefaultPages(&opt.Pages)
	internal.SetDefaultDomain(&opt.Domain)
	internal.SetDefaultStartPage(&opt.StartPage)
	internal.SetDefaultUserAgent(&opt.UserAgent)

	// Check validity of parameters.
	err := opt.checkParameterValidity()
	if err != nil {
		return nil, err
	}

	// Prepare payload with common parameters.
	payload := map[string]interface{}{
		"source":           oxylabs.GoogleShoppingPricing,
		"domain":           opt.Domain,
		"query":            query,
		"start_page":       opt.StartPage,
		"pages":            opt.Pages,
		"locale":           opt.Locale,
		"results_language": opt.ResultsLanguage,
		"geo_location":     opt.GeoLocation,
		"user_agent_type":  opt.UserAgent,
		"render":           opt.Render,
		"callback_url":     opt.CallbackURL,
		"parse":            opt.Parse,
	}

	// Add custom parsing instructions to the payload if provided.
	customParserFlag := false
	if opt.ParseInstructions != nil {
		payload["parsing_instructions"] = &opt.ParseInstructions
		customParserFlag = true
	}

	// Marshal.
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Req.
	httpResp, err := c.C.Req(ctx, jsonPayload, "POST")
	if err != nil {
		return nil, err
	}

	// Unmarshal the http Response and get the response.
	resp, err := GetResp(httpResp, opt.Parse, customParserFlag)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
