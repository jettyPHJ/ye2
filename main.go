package main

import (
	_ "file_handling/impl/task"
	"file_handling/utils/gin"
	_ "file_handling/utils/mysql"
)

func main() {
	gin.Router.Run(":12005")
}
