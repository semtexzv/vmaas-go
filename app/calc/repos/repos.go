package repos

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/calc"
	"regexp"
	"strings"
	"time"
)

func ReposByRegex(c *cache.Cache, repoPat string) []string {
	repoPat = strings.Trim(repoPat, "^$")
	regex := regexp.MustCompile(repoPat)

	res := []string{}
	for name := range c.RepoLabel2Ids {
		if regex.MatchString(name) {
			res = append(res, name)
		}
	}
	return res
}

func FilterModifiedSince(t time.Time) calc.Filter {
	return func(c *cache.Cache, id string) bool {
		ids := c.RepoLabel2Ids[id]
		for _, id := range ids {
			det := c.RepoDetails[id]
			return true
		}
	}
}
