package brokerclient

import (
	"log"
	"minireipaz/pkg/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaClient interface {
	Produce(topic string, key []byte, value []byte) error
	Close()
}

type KafkaClientImpl struct {
	producer *kafka.Producer
}

func NewBrokerClient(kafkaConfig config.KafkaConfig) KafkaClient {
	configMap := CreateKafkaConfigMap(kafkaConfig)

	producer, err := kafka.NewProducer(&configMap)
	if err != nil {
		log.Panicf("ERROR | Cannot create Broker producer %v", err)
	}

	return &KafkaClientImpl{
		producer: producer,
	}
}

func (k *KafkaClientImpl) Produce(topic string, key []byte, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}

	deliveryChan := make(chan kafka.Event)
	err := k.producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return err
}

func (k *KafkaClientImpl) Close() {
	k.producer.Close()
}

func CreateKafkaConfigMap(kafkaConfig config.KafkaConfig) kafka.ConfigMap {
	configMap := make(kafka.ConfigMap)
	configMap["bootstrap.servers"] = kafkaConfig.GetServersURI()
	configMap["security.protocol"] = kafkaConfig.GetProtocol()
	configMap["sasl.mechanisms"] = kafkaConfig.GetMechanisms()

	username := kafkaConfig.GetUsername()
	if username != "" {
		configMap["sasl.username"] = username
	}

	password := kafkaConfig.GetPassword()
	if password != "" {
		configMap["sasl.password"] = password
	}

	timeout := kafkaConfig.GetTimeout()
	if timeout != "" {
		configMap["session.timeout.ms"] = timeout
	}

	return configMap
}
