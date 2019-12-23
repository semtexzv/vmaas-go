package webserver

import (
	"encoding/json"
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"log"
	"net/http"
)

// start server serving loaded data
func Run() {
	log.Println("Starting simple webserver.")
	http.HandleFunc("/Packagename2Id", createHandler(cache.C.Packagename2Id))
	http.HandleFunc("/Id2Packagename", createHandler(cache.C.Id2Packagename))
	http.HandleFunc("/Updates", createHandler(cache.C.Updates))
	http.HandleFunc("/UpdatesIndex", createHandler(cache.C.UpdatesIndex))
	http.HandleFunc("/Evr2Id", createHandler(cache.C.Evr2Id))
	http.HandleFunc("/Id2Evr", createHandler(cache.C.Id2Evr))
	http.HandleFunc("/Id2Arch", createHandler(cache.C.Id2Arch))
	http.HandleFunc("/Arch2Id", createHandler(cache.C.Arch2Id))
	http.HandleFunc("/ArchCompat", createHandler(cache.C.ArchCompat))
	http.HandleFunc("/PackageDetails", createHandler(cache.C.PackageDetails))
	http.HandleFunc("/Nevra2PkgId", createHandler(cache.C.Nevra2PkgId))
	http.HandleFunc("/RepoDetails", createHandler(cache.C.RepoDetails))
	http.HandleFunc("/RepoLabel2Ids", createHandler(cache.C.RepoLabel2Ids))
	http.HandleFunc("/ProductId2RepoIds", createHandler(cache.C.ProductId2RepoIds))
	http.HandleFunc("/PkgId2RepoIds", createHandler(cache.C.PkgId2RepoIds))
	http.HandleFunc("/ErrataId2Name", createHandler(cache.C.ErrataId2Name))
	http.HandleFunc("/PkgId2ErrataIds", createHandler(cache.C.PkgId2ErrataIds))
	http.HandleFunc("/ErrataId2RepoIds", createHandler(cache.C.ErrataId2RepoIds))
	http.HandleFunc("/CveDetail", createHandler(cache.C.CveDetail))
	http.HandleFunc("/PkgErrata2Module", createHandler(cache.C.PkgErrata2Module))
	http.HandleFunc("/ModuleName2Ids", createHandler(cache.C.ModuleName2Ids))
	http.HandleFunc("/DbChange", createHandler(cache.C.DbChange))
	http.HandleFunc("/ErrataDetail", createHandler(cache.C.ErrataDetail))
	http.HandleFunc("/SrcPkgId2PkgId", createHandler(cache.C.SrcPkgId2PkgId))
	http.HandleFunc("/String", createHandler(cache.C.String))

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
    	print(err)
	}
}

// create handler returning json serialized value
func createHandler(value interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(js)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
