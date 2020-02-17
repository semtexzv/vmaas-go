package cache

import (
	"database/sql"
	"fmt"
	"github.com/RedHatInsights/vmaas-go/app/database"
	"github.com/RedHatInsights/vmaas-go/app/utils"
	"time"
)

var (
	C *Cache
)

func LoadCache() *Cache {
	c := Cache{}
	c.Packagename2Id = loadStrIntMap("packagename", "id", "packagename", "Packagename2Id")
	c.Id2Packagename = loadIntStrMap("packagename", "id", "packagename", "Id2Packagename")
	c.Updates = loadKeysValsOrderedMap("updates", "name_id", "package_id",
		"package_order", "name_id,package_order", "Updates")
	c.UpdatesIndex = loadKeysValsOrderedMap2Level("updates_index", "name_id", "evr_id", "package_order",
		"package_order", "name_id,package_order", "UpdatesIndex")
	c.Id2Evr, c.Evr2Id = loadEvrMaps("evr", "id", "Epoch", "Version", "Release",
		"Id2Evr, Evr2Id")
	c.Id2Arch = loadIntStrMap("arch", "id", "arch", "Id2Arch")
	c.Arch2Id = loadStrIntMap("arch", "id", "arch", "Arch2Id")
	c.ArchCompat = loadArchCompat()
	c.PackageDetails, c.Nevra2PkgId, c.SrcPkgId2PkgId = loadPkgDetails("PackageDetails, Nevra2PkgId, SrcPkgId2PkgId")
	c.RepoDetails, c.RepoLabel2Ids, c.ProductId2RepoIds = loadRepoDetails("RepoDetails, RepoLabel2Ids, ProductId2RepoIds")
	c.PkgId2RepoIds = loadInt2Ints("pkg_repo", "pkg_id,repo_id", "PkgId2RepoIds") // long
	c.ErrataDetail, c.ErrataId2Name = loadErratas("ErrataDetail, ErrataId2Name")
	c.PkgId2ErrataIds = loadInt2Ints("pkg_errata", "pkg_id,errata_id", "PkgId2ErrataIds")
	c.ErrataId2RepoIds = loadInt2Ints("errata_repo", "errata_id,repo_id", "ErrataId2RepoIds")
	c.CveDetail = loadCves("CveDetail")
	c.PkgErrata2Module = loadPkgErrataModule("PkgErrata2Module")
	c.ModuleName2Ids = loadModuleName2Ids("ModuleName2Ids")
	c.DbChange = loadDbChanges("DbChange")
	c.String = loadString("String")
	return &c
}

func getAllRows(tableName, cols, orderBy string) *sql.Rows {
	rows, err := database.Db.DB().Query(fmt.Sprintf("SELECT %s FROM %s ORDER BY %s",
		cols, tableName, orderBy))
	if err != nil {
		panic(err)
	}
	return rows
}

func loadIntArray(tableName, col, orderBy string) []int {
	rows := getAllRows(tableName, col, orderBy)
	defer rows.Close()

	var arr []int
	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			panic(err)
		}

		arr = append(arr, num)
	}
	return arr
}

func loadStrArray(tableName, col, orderBy string) []string {
	rows := getAllRows(tableName, col, orderBy)
	defer rows.Close()

	var arr []string
	for rows.Next() {
		var val string
		err := rows.Scan(&val)
		if err != nil {
			panic(err)
		}

		arr = append(arr, val)
	}
	return arr
}

func loadIntStrMap(tableName, intCol, strCol, info string) map[int]string {
	defer utils.TimeTrack(time.Now(), info)

	ints := loadIntArray(tableName, intCol, intCol)
	strs := loadStrArray(tableName, strCol, intCol)

	m := map[int]string{}

	for i := 0; i < len(ints); i++ {
		m[ints[i]] = strs[i]
	}
	return m
}

func loadStrIntMap(tableName, intCol, strCol, info string) map[string]int {
	defer utils.TimeTrack(time.Now(), info)

	ints := loadIntArray(tableName, intCol, strCol)
	strs := loadStrArray(tableName, strCol, strCol)

	m := map[string]int{}

	for i := 0; i < len(ints); i++ {
		m[strs[i]] = ints[i]
	}
	return m
}

func loadKeysValsOrderedMap(table, colKey, colVal, colOrder, orderBy, info string) map[int][]int {
	defer utils.TimeTrack(time.Now(), info)

	keys := loadIntArray(table, colKey, orderBy)
	vals := loadIntArray(table, colVal, orderBy)
	orders := loadIntArray(table, colOrder, orderBy)

	m := map[int][]int{}

	for i := 0; i < len(keys); i++ {
		if orders[i] == 0 {
			m[keys[i]] = make([]int, 0)
		}

		m[keys[i]] = append(m[keys[i]], vals[i])
	}
	return m
}

func loadKeysValsOrderedMap2Level(table, colKey, colNestkey, colVal, colOrder, orderBy, info string) map[int]map[int][]int {
	defer utils.TimeTrack(time.Now(), info)

	keys := loadIntArray(table, colKey, orderBy)
	nestKeys := loadIntArray(table, colNestkey, orderBy)
	vals := loadIntArray(table, colVal, orderBy)

	m := map[int]map[int][]int{}

	for i := 0; i < len(keys); i++ {
	   outer := m[keys[i]]
	   if outer == nil {
	      outer = map[int][]int{}
       }
       nested := outer[nestKeys[i]]
       nested = append(nested, vals[i])

       outer[nestKeys[i]] = nested
       m[keys[i]] = outer
	}
	return m
}

type Evr struct {
	Epoch   int
	Version string
	Release string
}

func loadEvrMaps(table, evrIdCol, epochCol, versionCol, releaseCol, info string) (map[int]Evr, map[Evr]int) {
	defer utils.TimeTrack(time.Now(), info)

	evrIds := loadIntArray(table, evrIdCol, evrIdCol)
	epochs := loadIntArray(table, epochCol, evrIdCol)
	vers := loadStrArray(table, versionCol, evrIdCol)
	rels := loadStrArray(table, releaseCol, evrIdCol)

	id2evr := map[int]Evr{}
	evr2id := map[Evr]int{}

	for i := 0; i < len(epochs); i++ {
		evr := Evr{
			Epoch:   epochs[i],
			Version: vers[i],
			Release: rels[i],
		}
		id2evr[evrIds[i]] = evr
		evr2id[evr] = evrIds[i]
	}
	return id2evr, evr2id
}

func loadArchCompat() map[int]map[int]bool {
	defer utils.TimeTrack(time.Now(), "arch_compat")

	orderBy := "from_arch_id,to_arch_id"
	fromArchIds := loadIntArray("arch_compat", "from_arch_id", orderBy)
	toArchIds := loadIntArray("arch_compat", "to_arch_id", orderBy)

	m := map[int]map[int]bool{}

	for i := 0; i < len(fromArchIds); i++ {
		from := m[fromArchIds[i]]
		if from == nil {
			from = map[int]bool{}
		}
		from[toArchIds[i]] = true
		m[fromArchIds[i]] = from
	}
	return m
}

func loadPkgDetails(info string) (map[int]PackageDetail, map[Nevra]int, map[int][]int) {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows("package_detail", "*", "id")
	id2pkdDetail := map[int]PackageDetail{}
	nevra2id := map[Nevra]int{}
	srcPkgId2PkgId := map[int][]int{}
	for rows.Next() {
		var pkgId int
		var det PackageDetail
		err := rows.Scan(&pkgId, &det.NameId, &det.EvrId, &det.ArchId, &det.SummaryId, &det.DescriptionId,
			&det.SrcPkgId)
		if err != nil {
			panic(err)
		}
		id2pkdDetail[pkgId] = det

		nevra := Nevra{det.NameId, det.EvrId, det.ArchId}
		nevra2id[nevra] = pkgId

		if det.SrcPkgId == nil {
			continue
		}

		_, ok := srcPkgId2PkgId[*det.SrcPkgId]
		if !ok {
			srcPkgId2PkgId[*det.SrcPkgId] = []int{}
		}

		srcPkgId2PkgId[*det.SrcPkgId] = append(srcPkgId2PkgId[*det.SrcPkgId], pkgId)
	}
	return id2pkdDetail, nevra2id, srcPkgId2PkgId
}

func loadRepoDetails(info string) (map[int]RepoDetail, map[string][]int, map[int][]int) {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows("repo_detail", "*", "label")
	id2repoDetail := map[int]RepoDetail{}
	repoLabel2id := map[string][]int{}
	prodId2RepoIds := map[int][]int{}
	for rows.Next() {
		var repoId int
		var det RepoDetail
		err := rows.Scan(&repoId, &det.Label, &det.Name, &det.Url, &det.BaseArch, &det.ReleaseVer,
			&det.Product, &det.ProductId, &det.Revision)
		if err != nil {
			panic(err)
		}
		id2repoDetail[repoId] = det

		_, ok := repoLabel2id[det.Label]
		if !ok {
			repoLabel2id[det.Label] = []int{}
		}
		repoLabel2id[det.Label] = append(repoLabel2id[det.Label], repoId)

		_, ok = prodId2RepoIds[det.ProductId]
		if !ok {
			prodId2RepoIds[det.ProductId] = []int{}
		}
		prodId2RepoIds[det.ProductId] = append(prodId2RepoIds[det.ProductId], repoId)
	}
	return id2repoDetail, repoLabel2id, prodId2RepoIds
}

func loadErratas(info string) (map[string]ErrataDetail, map[int]string) {
	defer utils.TimeTrack(time.Now(), info)

	erId2cves := loadInt2Strings("errata_cve", "errata_id,cve_id", "erId2cves")
	erId2pkgIds := loadInt2Ints("pkg_errata", "errata_id,pkg_id", "erId2pkgId")
	erId2modulePkgIds := loadInt2Ints("errata_modulepkg", "errata_id,pkg_id", "erId2modulePkgIds")
	erId2bzs := loadInt2Strings("errata_bugzilla", "errata_id,bugzilla", "erId2bzs")
	erId2refs := loadInt2Strings("errata_refs", "errata_id,ref", "erId2refs")
	erId2modules := loadErrataModules()

	cols := "id,name,synopsis,summary,type,severity,description,solution,issued,updated,url"
	rows := getAllRows("errata_detail", cols, "id")
	errataDetail := map[string]ErrataDetail{}
	errataId2Name := map[int]string{}
	for rows.Next() {
		var errataId int
		var errataName string
		var det ErrataDetail
		err := rows.Scan(&errataId, &errataName, &det.Synopsis, &det.Summary, &det.Type, &det.Severity,
			&det.Description, &det.Solution, &det.Issued, &det.Updated, &det.Url)
		if err != nil {
			panic(err)
		}
		errataId2Name[errataId] = errataName

		cves, ok := erId2cves[errataId]
		if ok {
			det.CVEs = cves
		}

		pkgIds, ok := erId2pkgIds[errataId]
		if ok {
			det.PkgIds = pkgIds
		}

		modulePkgIds, ok := erId2modulePkgIds[errataId]
		if ok {
			det.ModulePkgIds = modulePkgIds
		}

		bzs, ok := erId2bzs[errataId]
		if ok {
			det.Bugzillas = bzs
		}

		refs, ok := erId2refs[errataId]
		if ok {
			det.Refs = refs
		}

		modules, ok := erId2modules[errataId]
		if ok {
			det.Modules = modules
		}
		errataDetail[errataName] = det
	}
	return errataDetail, errataId2Name
}

func loadCves(info string) map[string]CveDetail {
	defer utils.TimeTrack(time.Now(), info)

	cveId2cwes := loadInt2Strings("cve_cwe", "cve_id,cwe", "cveId2cwes")
	cveId2pkg := loadInt2Ints("cve_pkg", "cve_id,pkg_id", "cveId2pkg")
	cve2eid := loadString2Ints("errata_cve", "cve_id,errata_id", "cve2eid")

	rows := getAllRows("cve_detail", "*", "id")
	cveDetails := map[string]CveDetail{}
	for rows.Next() {
		var cveId int
		var cveName string
		var det CveDetail
		err := rows.Scan(&cveId, &cveName, &det.RedHatUrl, &det.SecondaryUrl, &det.Cvss3Score, &det.Cvss3Metrics,
			&det.Impact, &det.PublishedDate, &det.ModifiedData, &det.Iava, &det.Description, &det.Cvss2Score,
			&det.Cvss2Metrics, &det.Source)
		if err != nil {
			panic(err)
		}

		cwes, ok := cveId2cwes[cveId]
		if ok {
			det.CWEs = cwes
		}

		pkgs, ok := cveId2pkg[cveId]
		if ok {
			det.PkgIds = pkgs
		}

		eids, ok := cve2eid[cveName]
		if ok {
			det.ErrataIds = eids
		}
		cveDetails[cveName] = det
	}
	return cveDetails
}

func loadPkgErrataModule(info string) map[PkgErrata][]int {
	defer utils.TimeTrack(time.Now(), info)

	orderBy := "pkg_id,errata_id,module_stream_id"
	table := "errata_modulepkg"
	pkgIds := loadIntArray(table, "pkg_id", orderBy)
	errataIds := loadIntArray(table, "errata_id", orderBy)
	moduleStreamIds := loadIntArray(table, "module_stream_id", orderBy)

	m := map[PkgErrata][]int{}

	for i := 0; i < len(pkgIds); i++ {
		pkgErrata := PkgErrata{pkgIds[i], errataIds[i]}
		_, ok := m[pkgErrata]
		if !ok {
			m[pkgErrata] = []int{}
		}

		m[pkgErrata] = append(m[pkgErrata], moduleStreamIds[i])
	}
	return m
}

func loadModuleName2Ids(info string) map[ModuleStream][]int {
	defer utils.TimeTrack(time.Now(), info)

	orderBy := "module,stream"
	table := "module_stream"
	modules := loadStrArray(table, "module", orderBy)
	streams := loadStrArray(table, "stream", orderBy)
	streamIds := loadIntArray(table, "stream_id", orderBy)

	m := map[ModuleStream][]int{}

	for i := 0; i < len(modules); i++ {
		pkgErrata := ModuleStream{modules[i], streams[i]}
		_, ok := m[pkgErrata]
		if !ok {
			m[pkgErrata] = []int{}
		}

		m[pkgErrata] = append(m[pkgErrata], streamIds[i])
	}
	return m
}

func loadString(info string) map[int]string {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows("string", "*", "id")
	m := map[int]string{}
	for rows.Next() {
		var id int
		var str *string
		err := rows.Scan(&id, &str)
		if err != nil {
			panic(err)
		}
		if str != nil {
           m[id] = *str
        }
	}
	return m
}

func loadDbChanges(info string) []DbChange {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows("dbchange", "*", "errata_changes")
	arr := []DbChange{}
	for rows.Next() {
		var item DbChange
		err := rows.Scan(&item.ErrataChanges, &item.CveChanges, &item.RepoChanges,
			&item.LastChange, &item.Exported)
		if err != nil {
			panic(err)
		}
		arr = append(arr, item)
	}
	return arr
}

func loadInt2Ints(table, cols, info string) map[int][]int {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows(table, cols, cols)
	int2ints := map[int][]int{}
	for rows.Next() {
		var key int
		var val int
		err := rows.Scan(&key, &val)
		if err != nil {
			panic(err)
		}

		_, ok := int2ints[key]
		if !ok {
			int2ints[key] = []int{}
		}
		int2ints[key] = append(int2ints[key], val)
	}
	return int2ints
}

func loadInt2Strings(table, cols, info string) map[int][]string {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows(table, cols, cols)
	int2strs := map[int][]string{}
	for rows.Next() {
		var key int
		var val string
		err := rows.Scan(&key, &val)
		if err != nil {
			panic(err)
		}

		_, ok := int2strs[key]
		if !ok {
			int2strs[key] = []string{}
		}

		int2strs[key] = append(int2strs[key], val)
	}
	return int2strs
}

func loadString2Ints(table, cols, info string) map[string][]int {
	defer utils.TimeTrack(time.Now(), info)

	rows := getAllRows(table, cols, cols)
	int2strs := map[string][]int{}
	for rows.Next() {
		var key string
		var val int
		err := rows.Scan(&key, &val)
		if err != nil {
			panic(err)
		}

		_, ok := int2strs[key]
		if !ok {
			int2strs[key] = []int{}
		}

		int2strs[key] = append(int2strs[key], val)
	}
	return int2strs
}

type Module struct {
	Name              string
	Stream            string
	Version           string
	Context           string
	PackageList       []string
	SourcePackageList []string
}

func loadErrataModules() map[int][]Module {
	defer utils.TimeTrack(time.Now(), "errata2module")

	rows := getAllRows("errata_module", "*", "errata_id")

	erId2modules := map[int][]Module{}
	for rows.Next() {
		var erId int
		var mod Module
		err := rows.Scan(&erId, &mod.Name, &mod.Stream, &mod.Version, &mod.Context)
		if err != nil {
			panic(err)
		}

		_, ok := erId2modules[erId]
		if !ok {
			erId2modules[erId] = []Module{}
		}

		erId2modules[erId] = append(erId2modules[erId], mod)
	}
	return erId2modules
}
