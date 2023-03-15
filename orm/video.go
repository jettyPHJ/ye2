package orm

// 定义视频结构体
type Video struct {
	ID         string `json:"id" gorm:"primaryKey"`
	UserID     string `json:"userId"`
	OrgaID     string `json:"OrgaId"`
	EngineFlag bool   `json:"engineFlag"`
	FileName   string `json:"fileName"`
	Path       string `json:"path"`
	Format     string `json:"format"`
	Size       string `json:"size"`
	Region     string `json:"region"`
	Section    string `json:"section"`
	Scenario   string `json:"scenario"`
	CreateTime int64  `json:"createTime"`
}
