package cache

import (
	"fmt"
)

func (c *Cache) Inspect() {
	fmt.Println("Packagename2Id:", len(c.Packagename2Id))
	fmt.Println("Id2Packagename:", len(c.Id2Packagename))
	if len(c.Packagename2Id) != len(c.Id2Packagename) {
		panic("not equal")
	}

	fmt.Println("Updates:", len(c.Updates))
	fmt.Println("UpdatesIndex:", len(c.UpdatesIndex))
	if len(c.Updates) != len(c.UpdatesIndex) {
		panic("not equal")
	}

	fmt.Println("Evr2Id:", len(c.Evr2Id))
	fmt.Println("Id2Evr:", len(c.Id2Evr))
	if len(c.Evr2Id) != len(c.Id2Evr) {
		panic("not equal")
	}

	fmt.Println("Arch2Id:", len(c.Arch2Id))
	fmt.Println("Id2Arch:", len(c.Id2Arch))
	if len(c.Arch2Id) != len(c.Id2Arch) {
		panic("not equal")
	}

	fmt.Println("ArchCompat:", len(c.ArchCompat))
	fmt.Println("PackageDetails:", len(c.PackageDetails))
	fmt.Println("Nevra2PkgId:", len(c.Nevra2PkgId))
	fmt.Println("RepoDetails:", len(c.RepoDetails))
	fmt.Println("RepoLabel2Ids:", len(c.RepoLabel2Ids))
	fmt.Println("ProductId2RepoIds:", len(c.ProductId2RepoIds))
	fmt.Println("PkgId2RepoIds:", len(c.PkgId2RepoIds))
	fmt.Println("ErrataId2Name:", len(c.ErrataId2Name))
	fmt.Println("PkgId2ErrataIds:", len(c.PkgId2ErrataIds))
	fmt.Println("ErrataId2RepoIds:", len(c.ErrataId2RepoIds))
	fmt.Println("CveDetail:", len(c.CveDetail))
	fmt.Println("PkgErrata2Module:", len(c.PkgErrata2Module))
	fmt.Println("ModuleName2Ids:", len(c.ModuleName2Ids))
	fmt.Println("DbChange:", len(c.DbChange))
	fmt.Println("ErrataDetail:", len(c.ErrataDetail))
	fmt.Println("SrcPkgId2PkgId:", len(c.SrcPkgId2PkgId))
	fmt.Println("String:", len(c.String))
}
