package task_release

import (
	"file_handling/orm"
	mysql "file_handling/utils/mysql"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type DBInserter struct{}

//新增video元信息记录和默认任务
func (*DBInserter) InsertDefault(videoMeta orm.Video) error {
	//开启事务处理
	tx := mysql.DataBase.Begin()
	//新增video元信息记录
	videoMeta.CreateTime = time.Now().Unix()
	result := tx.Create(videoMeta)
	if result.Error != nil {
		tx.Rollback()
	}
	//新增任务
	newTask := &orm.Task{
		ID:                 uuid.New().String(),
		VideoID:            videoMeta.ID,
		VideoPath:          videoMeta.Path,
		CurrentStage:       1,
		Status:             1,
		CreateTime:         time.Now().Unix(),
		TransportFilePath:  ChangeFileExt(videoMeta.Path, "_1.mp4"),
		TrajectoryFilePath: ChangeFileExt(videoMeta.Path, "_1.csv"),
		BackgroundFilePath: ChangeFileExt(videoMeta.Path, "_1.jpg"),
	}
	result = tx.Create(newTask)
	if result.Error != nil {
		tx.Rollback()
	}
	// 提交事务
	return tx.Commit().Error
}

//更改路径字符串中最后的文件名

func ChangeFileExt(path string, newExt string) string {
	dir, filename := filepath.Split(path)
	extension := filepath.Ext(filename)
	newFilename := filename[0:len(filename)-len(extension)] + newExt
	newPath := filepath.Join(dir, newFilename)
	return newPath
}
