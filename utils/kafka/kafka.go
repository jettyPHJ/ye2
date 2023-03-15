package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var Producer sarama.SyncProducer
var Consumer sarama.Consumer

func init() {
	// 创建 Kafka 生产者配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	// 创建 Kafka 生产者
	var err error
	Producer, err = sarama.NewSyncProducer([]string{"192.168.1.163:9092"}, config)
	if err != nil {
		panic(err)
	}
	Consumer, err = sarama.NewConsumer([]string{"192.168.1.163:9092"}, config)
	if err != nil {
		panic(err)
	}
}

//使用kfk发布消息
func SendMessage(topic string, data []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	_, _, err := Producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

//使用kfk消费消息
func ConsumeMsg(topic string, msgQueue chan []byte) error {
	partitionList, err := Consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, partition := range partitionList {
		partitionConsumer, err := Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				//fmt.Printf("Received message from partition %d at offset %d: %s\n", msg.Partition, msg.Offset, string(msg.Value))
				msgQueue <- msg.Value
			case err := <-partitionConsumer.Errors():
				fmt.Printf("Error while consuming partition %d: %s\n", partition, err.Error())
				return err
			}
		}
	}
	return nil
}
