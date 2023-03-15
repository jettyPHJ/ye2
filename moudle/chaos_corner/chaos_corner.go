package chaoscorner

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

//如果出现错误立即返回前端，停止协程
func RespondWithError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		runtime.Goexit()
	}
}

func DeserializeFromReader[T any](reader io.Reader) (T, error) {
	b := &bytes.Buffer{}
	io.Copy(b, reader)
	return DeserializeFromJSON[T](b.Bytes())
}

func DeserializeFromJSON[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
