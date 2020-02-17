package cve

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/calc"
	"regexp"
	"strings"
	"time"
)

func FilterExists() calc.Filter {
	return func(c *cache.Cache, cve string) bool {
		_, has := c.CveDetail[cve]
		return has
	}
}

func FilterRhOnly() calc.Filter {
	return func(c *cache.Cache, cve string) bool {
		return c.CveDetail[cve].Source == "Red Hat"
	}
}

func FilterModifiedSince(t time.Time) calc.Filter {
	return func(c *cache.Cache, cve string) bool {
		det := c.CveDetail[cve]
		if det.ModifiedDate != nil {
			return det.ModifiedDate.After(t)
		}
		if det.PublishedDate != nil {
			return det.PublishedDate.After(t)
		}
		return false
	}
}

func FilterPublishedSince(t time.Time) calc.Filter {
	return func(c *cache.Cache, cve string) bool {
		det := c.CveDetail[cve]
		if det.PublishedDate != nil {
			return det.PublishedDate.After(t)
		}
		return false
	}
}

type Request struct {
	CveList        []string   `json:"cve_list"`
	ModifiedSince  *time.Time `json:"modified_since"`
	PublishedSince *time.Time `json:"published_since"`
	RhOnly         bool       `json:"rh_only"`
	calc.Paging
}

type Response struct {
	CveList        map[string]cache.CveDetail `json:"cve_list"`
	ModifiedSince  *time.Time                 `json:"modified_since,omitempty"`
	PublishedSince *time.Time                 `json:"published_since,omitempty"`
	calc.Paging
}

func CvesByRegex(c *cache.Cache, patternStr string) []string {
	patternStr = strings.Trim(patternStr, "^$")
	regex := regexp.MustCompile(patternStr)

	res := []string{}
	for name := range c.CveDetail {
		if regex.MatchString(name) {
			res = append(res, name)
		}
	}
	return res
}

func Cves(c *cache.Cache, request Request) (Response, error) {
	resp := Response{}

	if len(request.CveList) == 0 {
		return resp, nil
	}

	if len(request.CveList) == 1 {
		request.CveList = CvesByRegex(c, request.CveList[0])
	}

	filters := []calc.Filter{FilterExists()}

	if request.RhOnly {
		filters = append(filters, FilterRhOnly())
	}

	if request.ModifiedSince != nil {
		filters = append(filters, FilterModifiedSince(*request.ModifiedSince))
		resp.ModifiedSince = request.ModifiedSince
	}

	if request.PublishedSince != nil {
		filters = append(filters, FilterPublishedSince(*request.PublishedSince))
		resp.PublishedSince = request.PublishedSince
	}

	cvePage, pagination := calc.Paginate(c, request.CveList, filters)
	for _, cve := range cvePage {
		det := c.CveDetail[cve]
		resp.CveList[cve] = det
	}
	resp.Paging = pagination
	return resp, nil
}
