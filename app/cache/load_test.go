package cache

import (
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/RedHatInsights/vmaas-go/app/database"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadCache(t *testing.T) {
	config.SQLiteFilePath = "../../data/vmaas.db"
	database.Configure()
	c := LoadCache()
	assert.Equal(t, 30173, len(c.Id2Packagename))
	assert.Equal(t, 30173, len(c.Packagename2Id))
	assert.Equal(t, 30173, len(c.Updates))
	assert.Equal(t, 30173, len(c.UpdatesIndex))
	assert.Equal(t, 67952, len(c.Evr2Id))
	assert.Equal(t, 67952, len(c.Id2Evr))
	assert.Equal(t, 33, len(c.Id2Arch))
	assert.Equal(t, 33, len(c.Arch2Id))
	assert.Equal(t, 28, len(c.ArchCompat))
	assert.Equal(t, 761595, len(c.PackageDetails))
	assert.Equal(t, 77468, len(c.SrcPkgId2PkgId))
	assert.Equal(t, 761595, len(c.Nevra2PkgId))
	assert.Equal(t, 20385, len(c.RepoDetails))
	assert.Equal(t, 3860, len(c.RepoLabel2Ids))
	assert.Equal(t, 332, len(c.ProductId2RepoIds))
	assert.Equal(t, 684024, len(c.PkgId2RepoIds)) // long
	assert.Equal(t, 26939, len(c.ErrataDetail))
	assert.Equal(t, 26939, len(c.ErrataId2Name))
	assert.Equal(t, 624278, len(c.PkgId2ErrataIds))
	assert.Equal(t, 26939, len(c.ErrataId2RepoIds))
	assert.Equal(t, 133911, len(c.CveDetail))
	assert.Equal(t, 7601, len(c.PkgErrata2Module))
	assert.Equal(t, 63, len(c.ModuleName2Ids))
	assert.Equal(t, 1, len(c.DbChange))
	assert.Equal(t, 50430, len(c.String))
}
