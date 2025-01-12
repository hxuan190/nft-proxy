package ginhttp

import (
	"github.com/alphabatem/nft-proxy/share"
	"net/http"
)

func (svc *httpService) paramErr(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, share.ResponseData(err))

}
