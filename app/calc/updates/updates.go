package updates

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/utils"
	"strings"
)

type Request struct {
	Packages []string             `json:"package_list"`
	RepoList []string             `json:"repository_list"`
	Modules  []cache.ModuleStream `json:"modules_list"`

	Releasever *string `json:"releasever"`
	BaseArch   *string `json:"basearch"`
}

type Update struct {
	Package    string `json:"package"`
	Erratum    string `json:"erratum"`
	Repository string `json:"repository"`
	Basearch   string `json:"basearch"`
	Releasever string `json:"releasever"`
}

type NameUpdateDetail struct {
	AvailableUpdates *[]Update `json:"available_updates,omitempty"`

	Description string `json:"description,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

type Response struct {
	UpdateList map[string]NameUpdateDetail `json:"update_list"`
	RepoList   []string                    `json:"repository_list"`
	ModuleList []cache.ModuleStream        `json:"modules_list,omitempty"`
	Releasever *string                     `json:"releasever,omitempty"`
	BaseArch   *string                     `json:"basearch,omitempty"`
}

func nil2empty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func nil2debug(s *string) string {
	if s != nil {
		return *s
	}
	return "<NIL>"
}

func ProcessRepositories(c *cache.Cache, req Request, resp *Response) (map[int]bool, error) {
	repos := make(map[int]bool)

	if len(req.RepoList) > 0 {
		for _, r := range req.RepoList {
			if ids, has := c.RepoLabel2Ids[r]; has {
				for _, l := range ids {
					repos[l] = true
				}
			}
		}
		resp.RepoList = req.RepoList
	} else {
		for id := range c.RepoDetails {
			repos[id] = true
		}
	}

	// Slow but works, iterate over every item, ANDing it with another condition
	// WARNING ! do string comparison instead of ptr comparison
	if req.Releasever != nil {
		for r := range repos {
			det := c.RepoDetails[r]
			repos[r] = repos[r] && det.ReleaseVer == nil && strings.Contains(det.Url, *req.Releasever) ||
				*det.ReleaseVer == *req.Releasever
		}
		resp.Releasever = req.Releasever
	}

	if req.BaseArch != nil {
		for r := range repos {
			det := c.RepoDetails[r]
			repos[r] = repos[r] && det.BaseArch == nil && strings.Contains(det.Url, *req.BaseArch) ||
				*det.BaseArch == *req.BaseArch

		}
		resp.BaseArch = req.BaseArch
	}

	for r, keep := range repos {
		if !keep {
			delete(repos, r)
		}
	}
	return repos, nil
}

func ProcessInputPackages(cache *cache.Cache, data Request, resp *Response) map[string]utils.Nevra {
	pkgs := map[string]utils.Nevra{}

	for _, p := range data.Packages {
		nevra, err := utils.ParseNevra(p)
		if resp.UpdateList == nil {
			resp.UpdateList = map[string]NameUpdateDetail{}
		}
		resp.UpdateList[p] = NameUpdateDetail{}
		if err != nil {
			continue
		}
		if _, has := cache.Packagename2Id[nevra.Name]; !has {
			continue
		}
		pkgs[p] = *nevra
	}
	return pkgs
}

func getRelatedProducts(c *cache.Cache, repoids map[int]bool) map[int]bool {
	products := make(map[int]bool)
	for p := range repoids {
		products[c.RepoDetails[p].ProductId] = true
	}
	return products
}

func getReleaseVersions(c *cache.Cache, originalRepoIds map[int]bool) map[string]bool {
	releasevers := map[string]bool{}
	for rid := range originalRepoIds {
		rv := c.RepoDetails[rid].ReleaseVer
		releasevers[nil2empty(rv)] = true
	}
	return releasevers
}

func getRepositories(c *cache.Cache, updatePkgid int, productIds map[int]bool, relesevers map[string]bool, errataIds []int, availableRepoIds map[int]bool) map[int]bool {
	repoids := map[int]bool{}
	errataRepoIds := map[int]bool{}

	for _, errataId := range errataIds {
		for _, rid := range c.ErrataId2RepoIds[errataId] {
			if availableRepoIds[rid] {
				errataRepoIds[rid] = true
			}
		}
	}

	for _, rid := range c.PkgId2RepoIds[updatePkgid] {
		if errataRepoIds[rid] {
			detail := c.RepoDetails[rid]
			if productIds[detail.ProductId] && relesevers[nil2empty(detail.ReleaseVer)] {
				repoids[rid] = true
			}
		}
	}
	return repoids
}

func buildNevra(c *cache.Cache, updatePkgId int) string {
	det := c.PackageDetails[updatePkgId]
	name := c.Id2Packagename[det.NameId]
	evr := c.Id2Evr[det.EvrId]
	arch := c.Id2Arch[det.ArchId]
	nevra := utils.Nevra{
		Name:    name,
		Epoch:   evr.Epoch,
		Version: evr.Version,
		Release: evr.Release,
		Arch:    arch,
	}
	return nevra.String()
}

func checkSecurityOnly(c *cache.Cache, securityOnly bool, errataId int) bool {
	if !securityOnly {
		return true
	}
	errataName := c.ErrataId2Name[errataId]
	if c.ErrataDetail[errataName].Type == "security" || len(c.ErrataDetail[errataName].CVEs) != 0 {
		return true
	}

	return false
}
func checkModules(c *cache.Cache, modules map[int]bool, updatePkgId int, errataId int) bool {

	if len(modules) == 0 {
		return true
	}

	pkgErrata := cache.PkgErrata{
		PkgId:    int(updatePkgId),
		ErrataId: int(errataId),
	}
	mods, has := c.PkgErrata2Module[pkgErrata]
	if !has {
		// Not a module related errata-pkg pair
		return true
	}

	for _, m := range mods {
		if has := modules[m]; !has {
			return false
		}
	}
	return true
}

func ProcessUpdates(c *cache.Cache,
	pkgs map[string]utils.Nevra,
	repos map[int]bool,
	modules map[int]bool,
	response *Response,
	includeTexts,
	securityOnly bool,
) error {
	for pkgString, pkg := range pkgs {
		evr := cache.Evr{
			Epoch:   pkg.Epoch,
			Version: pkg.Version,
			Release: pkg.Release,
		}

		nameId := c.Packagename2Id[pkg.Name]
		evrId := c.Evr2Id[evr]
		archId := c.Arch2Id[pkg.Arch]

		currentEvrIndexes := c.UpdatesIndex[nameId][evrId]
		if len(currentEvrIndexes) == 0 {
			continue
		}

		currentNevraPkgId := int(0)
		for _, idx := range currentEvrIndexes {
			pkgId := c.Updates[nameId][idx]
			// The package with same arch as searched from list specified in updatesIndex
			if c.PackageDetails[pkgId].ArchId == archId {
				currentNevraPkgId = pkgId
			}
		}

		if currentNevraPkgId == 0 {
			continue
		}

		response.UpdateList[pkgString] = NameUpdateDetail{
			AvailableUpdates: &[]Update{},
		}

		if includeTexts {
			updateList := response.UpdateList[pkgString]
			updateList.Summary = c.String[c.PackageDetails[currentNevraPkgId].SummaryId]
			updateList.Description = c.String[c.PackageDetails[currentNevraPkgId].DescriptionId]
			response.UpdateList[pkgString] = updateList
		}

		lastVersionPkgId := c.Updates[nameId][len(c.Updates[nameId])-1]
		if lastVersionPkgId == currentNevraPkgId {
			continue
		}

		originalPkgRepoids := map[int]bool{}
		for _, id := range c.PkgId2RepoIds[currentNevraPkgId] {
			originalPkgRepoids[id] = true
		}

		productIds := getRelatedProducts(c, originalPkgRepoids)
		validReleaseVers := getReleaseVersions(c, originalPkgRepoids)

		updatePkgIds := c.Updates[nameId][currentEvrIndexes[len(currentEvrIndexes)-1]+1:]

		for _, updatePkgId := range updatePkgIds {
			if _, has := c.PkgId2ErrataIds[updatePkgId]; !has {
				continue
			}
			updatedNevraArchId := c.PackageDetails[updatePkgId].ArchId

			if c.ArchCompat[archId][updatedNevraArchId] {
				continue
			}

			errataIds := c.PkgId2ErrataIds[updatePkgId]
			pkgNevra := buildNevra(c, updatePkgId)

			for _, errataId := range errataIds {
				if !checkSecurityOnly(c, securityOnly, errataId) {
					continue
				}
				if !checkModules(c, modules, updatePkgId, errataId) {
					continue
				}
				repoIds := getRepositories(c, updatePkgId, productIds, validReleaseVers, []int{errataId}, repos)

				for r := range repoIds {
					detail := c.RepoDetails[r]
					updates := response.UpdateList[pkgString]
					tmp := append(*updates.AvailableUpdates, Update{
						Package:    pkgNevra,
						Erratum:    c.ErrataId2Name[errataId],
						Repository: detail.Label,

						Basearch:   nil2empty(detail.BaseArch),
						Releasever: nil2empty(detail.ReleaseVer),
					})
					updates.AvailableUpdates = &tmp
					response.UpdateList[pkgString] = updates
				}

			}
		}

	}
	return nil
}

func Updates(cache *cache.Cache, data Request) (Response, error) {
	response := Response{}
	pkgs := ProcessInputPackages(cache, data, &response)

	repos, err := ProcessRepositories(cache, data, &response)

	modules := map[int]bool{}
	if len(data.Modules) > 0 {
		for _, m := range data.Modules {
			for _, mid := range cache.ModuleName2Ids[m] {
				modules[mid] = true
			}
		}
		response.ModuleList = data.Modules
	}

	err = ProcessUpdates(cache, pkgs, repos, modules, &response, false, false)
	return response, err
}