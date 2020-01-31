package calc

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/utils"
)

type UpdateReq struct {
	RepoList *[]string
	Packages []string
}

type Update struct {
	Package    string
	Erratum    string
	Repository string
	Basearch   string
	Releasever string
}
type NameUpdateDetail struct {
	AvailableUpdates []Update
}

type UpdatesRes struct {
	UpdateList map[string]NameUpdateDetail
}

func ProcessRepositories(cache *cache.Cache, req UpdateReq) ([]int, error) {
	res := []int{}
	if req.RepoList != nil && len(*req.RepoList) > 0 {
		for _, r := range *req.RepoList {
			if labels, has := cache.RepoLabel2Ids[r]; has {
				res = append(res, labels...)
			}

		}
	} else {
		for id := range cache.RepoDetails {
			res = append(res, id)
		}
	}
	return res, nil
}

func ProcessInputPackages(cache *cache.Cache, data UpdateReq) (map[string]utils.Nevra, error) {
	res := map[string]utils.Nevra{}

	for _, p := range data.Packages {
		nevra, err := utils.ParseNevra(p)
		if err != nil {
			return nil, err
		}
		if _, has := cache.Packagename2Id[nevra.Name]; !has {
			continue
		}
		res[p] = *nevra
	}
	return res, nil
}

func getRepositories(c *cache.Cache, updatePkgid int, productIds, errataIds, availableRepoIds []int) map[int]bool {
	repoids := map[int]bool{}
	errataRepoIds := map[int]bool{}

	for _, errataId := range errataIds {
		for _, rid := range c.ErrataId2RepoIds[errataId] {
			errataRepoIds[rid] = true
		}
	}

	for _, rid := range c.PkgId2RepoIds[updatePkgid] {
		if errataRepoIds[rid] {
			//detail := c.RepoDetails[rid]
			// TODO: Filter releasevers
			// TODO: Filter productid
			repoids[rid] = true
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
func ProcessUpdates(c *cache.Cache, repos []int, pkgs map[string]utils.Nevra) (UpdatesRes, error) {
	res := UpdatesRes{UpdateList: map[string]NameUpdateDetail{}}

	for name, pkg := range pkgs {
		evr := cache.Evr{
			Epoch:   pkg.Epoch,
			Version: pkg.Version,
			Release: pkg.Release,
		}
		nameId := c.Packagename2Id[pkg.Name]
		_ = c.Evr2Id[evr]
		archId := c.Arch2Id[pkg.Arch]
		// TODO: Missing secondary lookup by evr_id, investigate
		currentEvrIndexes := c.UpdatesIndex[nameId]
		if currentEvrIndexes == nil || len(currentEvrIndexes) == 0 {
			continue
		}

		currentNevraPkgId := 0
		for currentEvrIdx := range currentEvrIndexes {
			pkgId := c.Updates[nameId][currentEvrIdx]
			currNevraArchId := c.PackageDetails[pkgId].ArchId
			if currNevraArchId == archId {
				currentNevraPkgId = pkgId
				break
			}
		}

		if currentNevraPkgId == 0 {
			continue
		}

		lastVersionPkgId := c.Updates[nameId][len(c.Updates[nameId])-1 ]
		if lastVersionPkgId == currentNevraPkgId {
			continue
		}

		originalPkgRepoids := map[int]bool{}
		for _, id := range c.PkgId2RepoIds[currentNevraPkgId] {
			originalPkgRepoids[id] = true
		}

		//productIds := getRelatedProducts(originalPkgRepoids)
		//validReleaseVers := getRelatedProducts(originalPkgRepoids)

		updatePkgIds := c.Updates[nameId][currentEvrIndexes[len(currentEvrIndexes)-1]+1:]

		for _, updatePkgId := range updatePkgIds {
			if _, has := c.PkgId2ErrataIds[updatePkgId]; !has {
				continue
			}
			updatedNevraArchId := c.PackageDetails[updatePkgId].ArchId

			// TODO: Arch compat
			if archId != updatedNevraArchId {
				continue
			}

			errataIds := c.PkgId2ErrataIds[updatePkgId]
			pkgNevra := buildNevra(c, updatePkgId)
			for _, errataId := range errataIds {
				for r := range getRepositories(c, updatePkgId, []int{}, []int{errataId}, repos) {
					detail := c.RepoDetails[r]
					updates := res.UpdateList[name]
					updates.AvailableUpdates = append(updates.AvailableUpdates, Update{
						Package:    pkgNevra,
						Erratum:    c.ErrataId2Name[errataId],
						Repository: detail.Label,

						// TODO: Nil2empty
						Basearch:   "",
						Releasever: "",
					})
					res.UpdateList[name] = updates
				}

			}
		}

	}
	return res, nil
}

func CalcUpdates(cache *cache.Cache, data UpdateReq) (UpdatesRes, error) {
	pkgs, err := ProcessInputPackages(cache, data)

	if err != nil {
		panic(err)
	}
	repos, err := ProcessRepositories(cache, data)
	return ProcessUpdates(cache, repos, pkgs)
}
