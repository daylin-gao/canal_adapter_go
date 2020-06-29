package sysinit

import (
	"github.com/Shopify/sarama"
	"github.com/gao111/canal-adapter-go/config"
)

var (
	KafkaProducer sarama.SyncProducer
	KafkaConsumer sarama.Consumer
)

//init函数每个golang源文件中都可以定义一个init函数。golang系统中，所有的源文件都有自己所属的目录，每一个目录都有对应的包名。在包的引用中，一旦某一个包被使用，则这个包下边的init函数将会被执行，且只执行一次。只执行一次的含义是什么呢？如果一个包被多个地方引用，那么只有在这个包第一次被引用时，才会执行这个包里边的init函数，其他地方对包的再次引用，这个包里边的init函数不会被执行
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

	//使用配置,新建一个生产者
	KafkaProducer, err = sarama.NewSyncProducer([]string{config.Config.Kafka.Host+":"+config.Config.Kafka.Port}, kafkaConfig)
	if err != nil {
		panic(err)
	}

	KafkaConsumer, err = sarama.NewConsumer([]string{config.Config.Kafka.Host+":"+config.Config.Kafka.Port}, kafkaConfig)
	if err != nil {
		panic("error get consumer")
	}
	// 不关闭
	//defer KafkaProducer.Close()
}
