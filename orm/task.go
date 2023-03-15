package orm

// 定义任务结构体
type Task struct {
	ID                 string
	VideoID            string
	TransportFacID     string
	TransportDoneTime  int64
	TransportVideoID   string
	TrajectoryFacID    string
	TrajectoryDoneTime int64
	TrajectoryFileID   string
	BackgroundFacID    string
	BackgroundDoneTime int64
	BackgroundFileID   string
	CurrentStage       int
	Status             int
	CurrentProgress    string
	CreateTime         int64
	VideoPath          string
	TransportFilePath  string
	TrajectoryFilePath string
	BackgroundFilePath string
	ErrorCode          string
}
