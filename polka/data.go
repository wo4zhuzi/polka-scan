package polka

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"polka-scan/mongodb"
	"time"
)

func RunScan(concurrent int) error {

	height, err := GetLatestHeight()

	if err != nil {
		log.Fatal(err)
	}

	onlineHeight, err := GetOnlineHeight()

	if err != nil {
		log.Fatal(err)
	}

	//往前推两个块，防止块数据不全
	if height + concurrent > (onlineHeight - 2) {
		log.Printf("%s 未出块 \r\n", height + concurrent)
		time.Sleep(time.Duration(1)*time.Second)
		return nil
	}

	ch := make(chan mongodb.Aggregation, concurrent)
	for i:=height; i< height + concurrent ;i++  {
		//time.Sleep(time.Millisecond * time.Duration(100))
		go PushHeightData(i, ch)
	}

	var extrinsicList []mongodb.Extrinsics
	var eventsList []mongodb.Events
	var addressBalanceChange []mongodb.AddressBalanceChange

	for i:=height; i< height + concurrent ;i++  {
		aggregation := <-ch
		extrinsicList = append(extrinsicList, aggregation.ExtrinsicsList...)
		eventsList = append(eventsList, aggregation.EventsList...)
		addressBalanceChange = append(addressBalanceChange, aggregation.AddressBalanceChange...)
	}

	if len(extrinsicList) != 0 {
		//insertExtrinsicsList(extrinsicList)
	} else {
		log.Printf("高度 %d 获取数据异常、请检查节点", height);
		return err
	}

	if len(eventsList) != 0 {
		//insertEventsList(eventsList)
	}

	if len(addressBalanceChange) != 0 {
		insertAddressBalanceChange(addressBalanceChange)
	}

	upErr := UpdateBlockHeight(concurrent)

	if upErr == nil {
		log.Printf("height  %d ~ %d 插入成功\r\n", height, height + concurrent - 1)
	}

	return err
}

func PushHeightData(height int,ch chan mongodb.Aggregation)  {
	
	var extrinsicList []mongodb.Extrinsics
	var eventsList []mongodb.Events
	var addressBalanceChange []mongodb.AddressBalanceChange

	block, err := GetBlockData(height)

	if err != nil {
		log.Fatalf("高度 %d 获取数据失败，请检查节点\r\n", height)
	}

	for _, extrinsic := range block.Extrinsics  {
		extrinsicList = append(extrinsicList, extrinsic.ToMongoExtrinsics(block))
		eventsList = append(eventsList, extrinsic.ToMongoEvensList()...)
		addressBalanceChange = append(addressBalanceChange, extrinsic.ToMongoAddressBalanceChangeList(block)...)
	}

	ch <- mongodb.Aggregation{
		ExtrinsicsList:       extrinsicList,
		EventsList:           eventsList,
		AddressBalanceChange: addressBalanceChange,
	}

}

func GetLatestHeight() (height int, err error)  {
	collection := mongodb.MongoClient.Database("polkadot").Collection("block_height")

	var result mongodb.BlockHeight
	filter := bson.D{}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.Height, err
}

func insertExtrinsicsList(extrinsicList []mongodb.Extrinsics) ([]interface{}, error)  {
	collection := mongodb.MongoClient.Database("polkadot").Collection("extrinsics")

	var insertExtrinsicList []interface{}
	for _,v := range extrinsicList{
		insertExtrinsicList = append(insertExtrinsicList, v)
	}

	insertManyResult, err := collection.InsertMany(context.TODO(), insertExtrinsicList)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	return insertManyResult.InsertedIDs, err
}

func insertEventsList(eventsList []mongodb.Events) ([]interface{}, error)  {
	collection := mongodb.MongoClient.Database("polkadot").Collection("events")

	var insertEventsList []interface{}
	for _,v := range eventsList{
		insertEventsList = append(insertEventsList, v)
	}

	insertManyResult, err := collection.InsertMany(context.TODO(), insertEventsList)
	if err != nil {
		log.Fatal(err)
	}
	
	return insertManyResult.InsertedIDs, err
}


func insertAddressBalanceChange(addressBalanceChange []mongodb.AddressBalanceChange) ([]interface{}, error)  {
	collection := mongodb.MongoClient.Database("polkadot").Collection("address_balance_change")

	var insertAddressBalanceChange []interface{}
	for _,v := range addressBalanceChange{
		insertAddressBalanceChange = append(insertAddressBalanceChange, v)
	}


	insertManyResult, err := collection.InsertMany(context.TODO(), insertAddressBalanceChange)
	if err != nil {
		log.Fatal(err)
	}

	return insertManyResult.InsertedIDs, err
}

func UpdateBlockHeight(concurrent int) error {
	collection := mongodb.MongoClient.Database("polkadot").Collection("block_height")

	filter := bson.D{}

	update := bson.D{
		{"$inc", bson.D{
			{"height", concurrent},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return err
}

