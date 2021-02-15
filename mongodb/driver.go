package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"polka-scan/conf"
)

var MongoClient *mongo.Client

func init()  {
	mongoConf := conf.IniFile.Section("mongodb")
	ip := mongoConf.Key("ip").String()
	port := mongoConf.Key("port").String()

	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://"+ ip +":"+ port +"")

	var err error

	// 连接到MongoDB
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

}