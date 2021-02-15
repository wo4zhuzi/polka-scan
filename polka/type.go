package polka

const Decimal  = 10

//type Extrinsics struct {
//	ExtrinsicIndex     int    	`bson:"extrinsic_index"`
//	BlockNum           int    	`bson:"block_num"`
//	BlockTimestamp     int64  	`bson:"block_timestamp"`
//	ExtrinsicLength    string 	`bson:"extrinsic_length"`
//	VersionInfo        string 	`bson:"version_info"`
//	CallCode           string 	`bson:"call_code"`
//	CallModuleFunction string 	`bson:"call_module_function"`
//	CallModule         string 	`bson:"call_module"`
//	Params             string 	`bson:"params"`
//	AccountId          string 	`bson:"account_id"`
//	Signature          string 	`bson:"signature"`
//	Nonce              int	  	`bson:"nonce"`
//	Era                string 	`bson:"era"`
//	ExtrinsicHash      string 	`bson:"extrinsic_hash"`
//	IsSigned           int	  	`bson:"is_signed"`
//	Success            int	    `bson:"success"`
//	fee                float64	`bson:"fee"`
//}

type Block struct {
	Number			string
	Hash 			string
	ParentHash 		string
	StateRoot		string
	ExtrinsicsRoot  string
	AuthorId  		string
	Logs			[]Logs
	OnInitialize	interface{}
	Extrinsics		[]Extrinsics
	OnFinalize		interface{}
}

type Logs struct {
	Type 	string
	Index 	string
	Value	[]string
}

type Extrinsics struct {
	Method		Method
	Signature	interface{}
	Nonce		string
	Args		interface{}
	Tip			string
	Hash		string
	Info 		interface{}
	Events		[]Events
	Success		bool
	PaysFee 	bool
}

type Events struct {
	Method 	Method
	Data 	[]interface{}
}

type Method struct {
	Pallet 	string
	Method	string
}