package kafka

import (
	"fmt"
	"testing"
)

func TestKFK_producer(t *testing.T) {
	//Init()
	err := SendMessage("factoryFuncInfo", []byte("{sadfasdfasfsdddddddddddddddd}"))
	if err != nil {
		fmt.Println(err)
	}
}

func TestKFK_consumer(t *testing.T) {

}
