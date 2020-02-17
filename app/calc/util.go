package calc

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
)

type Paging struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Filter func(c *cache.Cache, id string) bool

func Paginate(c *cache.Cache, ids []string, filters []Filter) ([]string, Paging) {
	res := []string{}
Outer:
	for _, id := range ids {
		for _, f := range filters {
			if !f(c, id) {
				continue Outer
			}
		}
		res = append(res, id)
	}

	return res, Paging{}
}
