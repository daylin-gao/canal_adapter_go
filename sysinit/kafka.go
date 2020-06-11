package sysinit

import (
"errors"
"fmt"
	"github.com/jinzhu/gorm"
	"os"

"github.com/Shopify/sarama"
_ "github.com/jinzhu/gorm/dialects/mysql"
_ "github.com/jinzhu/gorm/dialects/postgres"
_ "github.com/mattn/go-sqlite3"
"github.com/gao111/canal-adapter-go/config"
)

var (
	KafkaProducer sarama.SyncProducer
	KafkaConsumer sarama.Consumer
)

func init() {
	var err error
	//设置配置
	kafkaConfig := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	kafkaConfig.Version = sarama.V0_11_0_0

	//使用配置,新建一个异步生产者
	KafkaProducer, err = sarama.NewSyncProducer([]string{config.Config.Kafka.Host+":"+config.Config.Kafka.Port}, kafkaConfig)
	if err != nil {
		panic(err)
	}
	// 不关闭
	//defer KafkaProducer.Close()
}
