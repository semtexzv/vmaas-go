package webserver

import (
	"github.com/RedHatInsights/vmaas-go/app/cache"
	"github.com/RedHatInsights/vmaas-go/app/calc/updates"
	"github.com/gin-gonic/gin"
	"time"
)

func up(c *gin.Context, req updates.Request) {
	res, err := updates.Updates(cache.C, req)
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
	} else {
		c.JSON(200, res)
	}
}

func Run() {
	r := gin.Default()

	r.GET("/api/v3/up/:nevra", func(c *gin.Context) {

		up(c, updates.Request{
			RepoList: []string{},
			Packages: []string{c.Param("nevra")},
		})

	})
	r.POST("/api/v3/up", func(c *gin.Context) {
		var req updates.Request
		err := c.BindJSON(&req)
		if err != nil {
			panic(req)
		}
		up(c, req)
	})
	go r.Run(":1080")
	<- time.NewTicker(20 * time.Second).C
}
