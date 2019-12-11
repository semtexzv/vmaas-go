package cache

import "time"

type Cache struct {
   Packagename2Id    map[string]int
   Id2Packagename    map[int]string
   Updates           map[int][]int
   UpdatesIndex      map[int][]int
   Evr2Id            map[Evr]int
   Id2Evr            map[int]Evr
   Id2Arch           map[int]string
   Arch2Id           map[string]int
   ArchCompat        map[int][]int
   PackageDetails    map[int]PackageDetail
   Nevra2PkgId       map[Nevra]int
   RepoDetails       map[int]RepoDetail
   RepoLabel2Ids     map[string][]int
   ProductId2RepoIds map[int][]int
   PkgId2RepoIds     map[int][]int
   ErrataId2Name     map[int]string
   PkgId2ErrataIds   map[int][]int
   ErrataId2RepoIds  map[int][]int
   CveDetail         map[string]CveDetail
   PkgErrata2Module  map[PkgErrata][]int
   ModuleName2Ids    map[ModuleStream][]int
   DbChange          []DbChange
   ErrataDetail      map[string]ErrataDetail
   SrcPkgId2PkgId    map[int][]int
   String            map[string]*string
}

type PackageDetail struct {
   NameId        int
   EvrId         int
   ArchId        int
   SummaryId     string
   DescriptionId string
   SrcPkgId      *int
}

type Nevra struct {
   NameId    int
   EvrId     int
   ArchId    int
}

type RepoDetail struct {
   Label      string
   Name       string
   Url        string
   BaseArch   *string
   ReleaseVer *string
   Product    string
   ProductId  int
   Revision   string
}

type CveDetail struct {
   RedHatUrl     *string
   SecondaryUrl  *string
   Cvss3Score    *float64
   Cvss3Metrics  *string
   Impact        string
   PublishedDate *time.Time
   ModifiedData  *time.Time
   Iava          *string
   Description   string
   Cvss2Score    *float64
   Cvss2Metrics  *string
   Source        string
   CWEs          []string
   PkgIds        []int
   ErrataIds     []int
}

type PkgErrata struct {
   PkgId    int
   ErrataId int
}

type ModuleStream struct {
   Module string
   Stream string
}

type DbChange struct {
   ErrataChanges time.Time
   CveChanges    time.Time
   RepoChanges   time.Time
   LastChange    time.Time
   Exported      time.Time
}

type ErrataDetail struct {
   Synopsis     string
   Summary      string
   Type         string
   Severity     string
   Description  *string
   CVEs         []string
   PkgIds       []int
   ModulePkgIds []int
   Bugzillas    []string
   Refs         []string
   Modules      []Module
   Solution     string
   Issued       time.Time
   Updated      time.Time
   Url          string
}
