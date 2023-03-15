package process

import (
	r "file_handling/utils/gin"

	"github.com/gin-gonic/gin"
)

func init() {
	v := r.Router.Group("/videoFile")
	{
		v.POST("", post)
	}
}

func post(c *gin.Context) {

}
