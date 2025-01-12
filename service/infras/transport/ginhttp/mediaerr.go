package ginhttp

import (
	"log"
	"net/http"
	"os"
)

// TODO Replace with placeholder image
func (svc *httpService) mediaError(c *gin.Context, err error) {
	log.Printf("Media Err: %s", err)

	img, err := os.ReadFile("./docs/failed_image.jpg")
	if err != nil {
		log.Println(err)
	}
	c.Header("Cache-Control", "public, max=age=60") //Stop flooding
	c.Data(http.StatusOK, "image/jpeg", img)
	//c.JSON(200, gin.H{
	//	"error": err.Error(),
	//})
}
