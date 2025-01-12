package ginhttp

import (
	"github.com/alphabatem/nft-proxy/share"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// // @Summary Ping liquify services
// // @Accept  json
// // @Produce json
// // @Router /nfts/{id} [get]
func (svc *httpService) ShowNFT(c *gin.Context) {
	svc.statHandler.IncrementMediaRequests()

	skipCache, _ := strconv.ParseBool(c.DefaultQuery("nocache", ""))
	if skipCache || rand.Intn(1000) == 1 {
		if err := svc.imageHandler.ClearCache(c.Param("id")); err != nil {
			svc.paramErr(c, err)
			return
		}
	}

	media, err := svc.imageHandler.Media(c.Param("id"), skipCache)
	if err != nil {
		svc.paramErr(c, err)
		return
	}

	c.Header("Cache-Control", "public, max-age=172800")
	c.Header("Expires", time.Now().AddDate(0, 0, 2).Format(http.TimeFormat))

	c.JSON(http.StatusOK, share.ResponseData(media))
}
