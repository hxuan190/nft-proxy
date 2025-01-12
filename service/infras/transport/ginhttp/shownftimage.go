package ginhttp

// // @Summary Ping liquify services
// // @Accept  json
// // @Produce json
// // @Router/nfts/{id}/image [get]
func (svc *httpService) ShowNFTImage(c *gin.Context) {
	svc.statHandler.IncrementImageFileRequests()
	err := svc.imageHandler.ImageFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}
