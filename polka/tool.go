package polka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"polka-scan/conf"
	"polka-scan/mongodb"
	"strconv"
)

func (p Extrinsics) ToMongoExtrinsics(block Block) mongodb.Extrinsics  {

	argsArr := p.Args.(map[string]interface{})
	destAddress  := ""

	var value float64 = 0
	_destAddress, exists := argsArr["dest"].(string)

	if exists {
		destAddress = _destAddress
	}
	_value, exists := argsArr["value"].(string)

	if exists {
		_value_float64, _ := strconv.ParseFloat(_value, 64)
		value = _value_float64 / math.Pow(10, Decimal)
	}

	//args, _ := json.Marshal(p.Args)

	signature, ok := p.Signature.(map[string]interface{})

	address := ""
	IsSigned := false
	if ok {
		address = signature["signer"].(string)
		IsSigned = true
	}

	info , ok := p.Info.(map[string]interface{})
	var fee int64 = 0
	if ok {
		partialFee, exists := info["partialFee"]

		if exists {
			fee, _ = strconv.ParseInt(partialFee.(string), 10, 64)
		}
	}

	BlockNum, _ := strconv.Atoi(block.Number)

	now, exists := block.Extrinsics[0].Args.(map[string]interface{})["now"]
	var BlockTimestamp int64 = 0
	if exists {
		_BlockTimestamp := now.(string)
		BlockTimestamp, _ = strconv.ParseInt(_BlockTimestamp[0:len(_BlockTimestamp) -3], 10, 64)
	}

	data := mongodb.Extrinsics{
		//ExtrinsicIndex:     0,
		BlockNum:           BlockNum,
		BlockTimestamp:     BlockTimestamp,
		//ExtrinsicLength:    "",
		//VersionInfo:        "",
		//CallCode:           "",
		CallModule:         p.Method.Pallet,
		CallModuleFunction: p.Method.Method,
		//Args:               string(args),
		Args:               "",
		Address:            address,
		//Signature:          "",
		Nonce:              p.Nonce,
		//Era:                "",
		ExtrinsicHash:      p.Hash,
		IsSigned:           IsSigned,
		Success:            p.Success,
		DestAddress:		destAddress,
		Value:				value,
		Fee:                fee,
	}

	return data
}

func (p Extrinsics) ToMongoAddressBalanceChangeList(block Block) []mongodb.AddressBalanceChange  {
	events := p.Events

	BlockNum, _ := strconv.Atoi(block.Number)

	now, exists := block.Extrinsics[0].Args.(map[string]interface{})["now"]
	var BlockTimestamp int64 = 0
	if exists {
		_BlockTimestamp := now.(string)
		BlockTimestamp, _ = strconv.ParseInt(_BlockTimestamp[0:len(_BlockTimestamp) -3], 10, 64)
	}

	info , ok := p.Info.(map[string]interface{})
	var fee int64 = 0
	if ok {
		partialFee, exists := info["partialFee"]

		if exists {
			fee, _ = strconv.ParseInt(partialFee.(string), 10, 64)
		}
	}

	var addressBalanceChangeList  []mongodb.AddressBalanceChange

	//处理交易失败的数据
	if p.Success == false {
		if p.Method.Pallet == "balances" && (p.Method.Method == "transferKeepAlive" || p.Method.Method == "transfer" || p.Method.Method == "forceTransfer") {

			signature := p.Signature.(map[string]interface{})
			args := p.Args.(map[string]interface{})

			_value_float64, _ := strconv.ParseFloat(args["value"].(string), 64)
			value := _value_float64 / math.Pow(10, Decimal)

			addressBalanceChange := mongodb.AddressBalanceChange{
				BlockNum:      BlockNum,
				Time:          BlockTimestamp,
				From:          signature["signer"].(string),
				To:            args["dest"].(string),
				Value:         value,
				Fee:           fee,
				ExtrinsicHash: p.Hash,
				Success:       false,
			}
			addressBalanceChangeList = append(addressBalanceChangeList, addressBalanceChange)
		}

	}

	for _, event := range events {
		if event.Method.Pallet == "balances" && event.Method.Method == "Transfer" {
			_value_float64, _ := strconv.ParseFloat(event.Data[2].(string), 64)
			value := _value_float64 / math.Pow(10, Decimal)
			addressBalanceChange := mongodb.AddressBalanceChange{
				BlockNum:      BlockNum,
				Time:          BlockTimestamp,
				From:          event.Data[0].(string),
				To:            event.Data[1].(string),
				Value:         value,
				Fee:           fee,
				ExtrinsicHash: p.Hash,
				Success: 	   p.Success,
			}
			addressBalanceChangeList = append(addressBalanceChangeList, addressBalanceChange)
		}
	}

	return addressBalanceChangeList
}

func (p Extrinsics) ToMongoEvensList() []mongodb.Events  {

	events := p.Events
	var eventsList []mongodb.Events

	for _, event := range events {
		data, _ := json.Marshal(event.Data)

		mongoEvents := mongodb.Events{
			ExtrinsicHash:      p.Hash,
			CallModule:        	p.Method.Pallet,
			CallModuleFunction: p.Method.Method,
			Data:               string(data),
		}

		eventsList = append(eventsList, mongoEvents)
	}

	return eventsList
}

func GetBlockData(blockNum int) (Block, error) {
	polkadotConf := conf.IniFile.Section("polkadot")
	api := polkadotConf.Key("api").String()

	resp, err := http.Get(""+ api +"/blocks/" + strconv.Itoa(blockNum))

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var block Block
	json.Unmarshal(body, &block)

	return block, err
}

func GetOnlineHeight() (int, error) {
	polkadotConf := conf.IniFile.Section("polkadot")
	api := polkadotConf.Key("api").String()

	resp, err := http.Get(""+ api +"/blocks/head")

	if err != nil {
		fmt.Println(err)
		return 0,err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var block Block
	json.Unmarshal(body, &block)

	_number := block.Number
	number, err := strconv.Atoi(_number)

	return number,err
}