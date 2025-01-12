package ginhttp

import (
	"github.com/alphabatem/nft-proxy/share"
	"net/http"
)

// // @Summary Ping liquify services
// // @Accept  json
// // @Produce json
// // @Router /stats [get]
func (svc *httpService) GetStat(c *gin.Context) {
	ctx := c.Request.Context()
	stats, err := svc.statHandler.ServiceStats(ctx)
	if err != nil {
		svc.paramErr(c, err)
		return
	}

	c.JSON(http.StatusOK, share.ResponseData(stats))
}
