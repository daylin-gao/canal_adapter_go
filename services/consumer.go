// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/gao111/canal-adapter-go/config"
	"github.com/gao111/canal-adapter-go/models"
	"log"
	"os"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/gao111/canal-adapter-go/sysinit"
	"github.com/gao111/canal-adapter-go/beans"
	protocol "github.com/gao111/canal-adapter-go/protocol"
	"github.com/golang/protobuf/proto"
)

func main() {
	//根据消费者获取指定的主题分区的消费者,Offset这里指定为获取最新的消息.
	partitionConsumer, err := sysinit.KafkaConsumer.ConsumePartition(config.Config.Kafka.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("error get partition consumer", err)
	}

	for {
		select {
		//接收消息通道和错误通道的内容.
		case msg := <-partitionConsumer.Messages():
			fmt.Println("msg offset: ", msg.Offset, " partition: ", msg.Partition, " timestrap: ", msg.Timestamp.Format("2006-Jan-02 15:04"), " value: ", string(msg.Value))
			var binlogBean beans.BinlogBean
			err := json.Unmarshal(msg.Value, &binlogBean)
			if err != nil {
				log.Printf("解析数据出错 , err:%s" , err)
			}

			go syncEntry(binlogBean.Columns , binlogBean.EventType , binlogBean.SchemaName , binlogBean.TableName)

		case err := <-partitionConsumer.Errors():
			fmt.Println(err.Err)
		}
	}
}


func syncEntry (columns []*protocol.Column ,eventType protocol.EventType ,dbName string , tableName string) {
	sync := models.NewSync(dbName , tableName)
	if eventType == protocol.EventType_DELETE {
		//fmt.Println(fmt.Sprintf("---------------------------%s : %s ", dbName, tableName))
		//printColumn(columns)
		sync.DeleteSync(columns)
	} else if eventType == protocol.EventType_INSERT {
		//printColumn(columns)
		sync.InsertSync(columns)
	} else if eventType == protocol.EventType_UPDATE {
		sync.UpdateSync(columns)
	}
}

func printEntry(entrys []protocol.Entry) {
	for _, entry := range entrys {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(protocol.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == protocol.EventType_DELETE {
					printColumn(rowData.GetBeforeColumns())
				} else if eventType == protocol.EventType_INSERT {
					printColumn(rowData.GetAfterColumns())
				} else {
					fmt.Println("-------> before")
					printColumn(rowData.GetBeforeColumns())
					fmt.Println("-------> after")
					printColumn(rowData.GetAfterColumns())
				}
			}
		}
	}
}

func printColumn(columns []*protocol.Column) {
	for _, col := range columns {
		fmt.Println(fmt.Sprintf("%s : %s  update= %t , mysql_type=%s", col.GetName(), col.GetValue(), col.GetUpdated(), col.GetSqlType()))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
