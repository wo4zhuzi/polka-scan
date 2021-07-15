package mongodb

//交易详情表
type Extrinsics struct {
	//ExtrinsicIndex     int    	`bson:"extrinsic_index"`
	BlockNum       int   `bson:"block_num"`
	BlockTimestamp int64 `bson:"block_timestamp"`
	//ExtrinsicLength    string 	`bson:"extrinsic_length"`
	//VersionInfo        string 	`bson:"version_info"`
	//CallCode           string 	`bson:"call_code"`
	CallModule         string `bson:"call_module"`
	CallModuleFunction string `bson:"call_module_function"`
	Args               string `bson:"args"`
	Address            string `bson:"address"`
	//Signature          string 	`bson:"signature"`
	Nonce string `bson:"nonce"`
	//Era                string 	`bson:"era"`
	ExtrinsicHash string  `bson:"extrinsic_hash"`
	IsSigned      bool    `bson:"is_signed"`
	Success       bool    `bson:"success"`
	DestAddress   string  `bson:"dest_address"`
	Value         float64 `bson:"amount"`
	Fee           int64   `bson:"fee"`
}

//账户金额变动表
type AddressBalanceChange struct {
	BlockNum      int     `bson:"block_num"`
	Time          int64   `bson:"time"`
	From          string  `bson:"from"`
	To            string  `bson:"to"`
	Value         float64 `bson:"value"`
	Fee           int64   `bson:"fee"`
	ExtrinsicHash string  `bson:"extrinsic_hash"`
	Success       bool    `bson:"success"`
}

//事件
type Events struct {
	ExtrinsicHash      string `bson:"extrinsic_hash"`
	CallModule         string `bson:"call_module"`
	CallModuleFunction string `bson:"call_module_function"`
	Data               string `bson:"data"`
}

//入库数据聚合
type Aggregation struct {
	ExtrinsicsList       []Extrinsics
	EventsList           []Events
	AddressBalanceChange []AddressBalanceChange
}

//最新高度表
type BlockHeight struct {
	Height int `bson:"height"`
}
