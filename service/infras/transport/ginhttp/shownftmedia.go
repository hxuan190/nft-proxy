package ginhttp

// // @Summary Ping liquify services
// // @Accept  json
// // @Produce json
// // @Router /nfts/{id}/media [get]
func (svc *httpService) ShowNFTMedia(c *gin.Context) {
	svc.statHandler.IncrementMediaFileRequests()
	err := svc.imageHandler.MediaFile(c, c.Param("id"))
	if err != nil {
		svc.mediaError(c, err)
		return
	}
}
