package orm

type FactoryInfoList struct {
	MsgID           string        `json:"msgID"`
	Total           int           `json:"total"`
	FactoryInfoList []FactoryInfo `json:"factoryInfoList"`
}

type FactoryInfo struct {
	FactoryID       string `json:"factoryID"`
	ParentFactoryID string `json:"parentFactoryID"`
	FactoryNumber   string `json:"factoryNumber"`
	FactoryType     string `json:"factoryType"`
	RoadNumber      string `json:"roadNumber"`
	FactoryName     string `json:"factoryName"`
	UpTime          string `json:"upTime"`
	CreatTime       string `json:"creatTime"`
	Factory         string `json:"factory"`
	Model           string `json:"model"`
	HardVersion     string `json:"hardVersion"`
	SoftVersion     string `json:"softVersion"`
	Temperature     string `json:"temperature"`
	PrivateData     string `json:"privateData"`
	ErrorCode       string `json:"errorCode"`
	Invalid         bool   `json:"invalid"`
}

type FuncTableRes struct {
	MsgID            string           `json:"msgID"`
	FactoryFuncTable FactoryFuncTable `json:"factoryFuncTable"`
}

type FactoryFuncTable struct {
	FactoryID string        `json:"factoryID"`
	FuncTable []FuncDetails `json:"funcTable"`
}

type FuncDetails struct {
	Function     string `json:"function"`
	Annotation   string `json:"annotation"`
	CurrentValue string `json:"currentValue"`
}