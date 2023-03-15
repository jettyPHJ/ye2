package orm

type RealeseTaskToEngine struct {
	FactoryID string `json:"factoryID"`
	Time      int64  `json:"time"`
	Function  string `json:"function"`
	Value     string `json:"value"`
}

type TaskDescription struct {
	MsgID        string `json:"msgID"`
	Time         int64  `json:"time"`
	SrcFilePath  string `json:"srcFilePath"`
	DestFilePath string `json:"destFilePath"`
}

type SubscriptionData struct {
	FactoryID string `json:"factoryID"`
	Data      string `json:"data"`
}

//SubscriptionData可能的数据类型：视频转码
type Data01 struct {
	MsgID    string `json:"msgID" required:"true"`
	CurFrame int    `json:"curFream" required:"true"`
	AllFrame int    `json:"allFream" required:"true"`
	Fps      int    `json:"fps" required:"true"`
	ErrCode  string `json:"errCode" required:"true"`
}

type UpdateProgress struct {
	TaskID   string
	Progress string
	DoneFlag bool
	Err      error
}

type TaskDescription02 struct {
	MsgID                 string `json:"msgID"`
	Time                  int64  `json:"time"`
	SrcVideoFilePath      string `json:"srcVideoFilePath"`
	SrcTrajectoryFilePath string `json:"srcTrajectoryFilePath"`
	DestImgFilePath       string `json:"destImgFilePath"`
}
