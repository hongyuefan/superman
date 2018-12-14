## Quick Start

#### 下载

    go get github.com/hongyuefan/superman

#### 搭建环境

本系统依赖kafka和influxDB，安装好后，kafka需要实现声明topic，使用下面的语句

    ./bin/kafka-topics.sh  --create  --zookeeper  localhost:2181  --replication-factor 1  --partitions  1  --topic okex_spider_data
	
    ./bin/kafka-topics.sh  --create  --zookeeper  localhost:2181  --replication-factor 1  --partitions  1  --topic okex_archer_req
	
    ./bin/kafka-topics.sh  --create  --zookeeper  localhost:2181  --replication-factor 1  --partitions  1  --topic okex_archer_rsp

influxDB安装下载见[官网](https://www.influxdata.com/)

#### 使用

进入build目录，使用build.sh生成可执行文件，如果生成失败，请自行安装所需要的第三方库：

	go get github.com/Shopify/sarama
	go get github.com/bitly/go-simplejson
	go get github.com/syndtr/goleveldb/leveldb
	go get github.com/gorilla/websocket

执行build/run.sh

## 代码说明

    archer  下单程序
    spider  订阅收集行情程序
    krang   运行策略和计算行情指标程序
    stg     行情存储，将交易所一天的行情全部存到一个leveldb数据库，这些数据用于回放
    replay  回放程序，用于调试策略
    strategy 策略模块，新加策略放到该模块下
