package enginejoint

import (
	"encoding/json"
	"errors"
	"file_handling/orm"
	myhttp "file_handling/utils/http"
	"file_handling/utils/kafka"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enginejoint struct {
	DB *gorm.DB
}

func (e *Enginejoint) SearchForTask(status int) []orm.Task {
	result := []orm.Task{}
	e.DB.Where("status = ? AND current_stage != ?", status, 4).Find(&result)
	return result
}

//返回所有符合task的factoryInfo
func (e *Enginejoint) SearchForFactoryID(task orm.Task) []orm.FactoryInfo {
	result := []orm.FactoryInfo{}
	pageNum := 1
	for {
		u, _ := url.Parse("http://" + "192.168.1.163:8000" + "/v1/factoryService/info")
		q := u.Query()
		q.Set("PageSize", "10")
		q.Set("PageNum", fmt.Sprint(pageNum))
		q.Set("MsgID", uuid.New().String())
		q.Set("FactoryType", "算法")
		switch task.CurrentStage {
		case 1:
			q.Set("FactoryName", "MP4编码生成器")
		case 2:
			q.Set("FactoryName", "轨迹文件生成算法")
		case 3:
			q.Set("FactoryName", "背景图片生成算法")
		}
		u.RawQuery = q.Encode()
		response, err := myhttp.HttpDo("GET", u.String(), nil)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		factoryInfoList, _ := myhttp.DeserializeFromHttpResponse[orm.FactoryInfoList](response)
		result = append(result, factoryInfoList.FactoryInfoList...)
		if len(result) == factoryInfoList.Total {
			return result
		}
		pageNum++
		response.Body.Close()
	}
}

//在factoryInfos中随机挑选一个可用的factory
func (e *Enginejoint) SelectOneFactory(factories []orm.FactoryInfo) (*orm.FactoryInfo, error) {
	available_fac := []orm.FactoryInfo{}
	for _, v := range factories {
		if !v.Invalid {
			available_fac = append(available_fac, v)
		}
	}
	if len(available_fac) == 0 {
		return nil, errors.New("没有可用的factory")
	}
	return &available_fac[rand.Intn(len(available_fac))], nil
}

//返回factory当前忙碌的任务ID-----数据引擎3.4.2接口
func (e *Enginejoint) CheckFactoryFree(factory orm.FactoryInfo, task orm.Task) (taskID string) {
	u, _ := url.Parse("http://" + "192.168.1.163:8000" + "/v1/factoryService/funcTable")
	q := u.Query()
	q.Set("MsgID", uuid.New().String())
	q.Set("FactoryID", factory.FactoryID)
	u.RawQuery = q.Encode()
	response, err := myhttp.HttpDo("GET", u.String(), nil)
	if err != nil {
		return ""
	}
	funcTable, _ := myhttp.DeserializeFromHttpResponse[orm.FuncTableRes](response)
	switch task.CurrentStage {
	case 1:
		for _, v := range funcTable.FactoryFuncTable.FuncTable {
			if v.Function == "生成MP4" {
				data := orm.Data01{}
				json.Unmarshal([]byte(v.CurrentValue), &data)
				return data.MsgID
			}
		}
	case 2:
		for _, v := range funcTable.FactoryFuncTable.FuncTable {
			if v.Function == "生成轨迹文件" {
				data := orm.Data01{}
				json.Unmarshal([]byte(v.CurrentValue), &data)
				return data.MsgID
			}
		}
	case 3:
		for _, v := range funcTable.FactoryFuncTable.FuncTable {
			if v.Function == "生成背景图片" {
				data := orm.Data01{}
				json.Unmarshal([]byte(v.CurrentValue), &data)
				return data.MsgID
			}
		}
	}
	return ""
}

//更新任务表:确认factoryId,将status置为2（ready）
func (e *Enginejoint) Update1(task orm.Task, fac orm.FactoryInfo) {
	switch task.CurrentStage {
	case 1:
		e.DB.Model(&task).Update("TransportFacID", fac.FactoryID)
	case 2:
		e.DB.Model(&task).Update("TrajectoryFacID", fac.FactoryID)
	case 3:
		e.DB.Model(&task).Update("BackgroundFacID", fac.FactoryID)
	}
	e.DB.Model(&task).Update("Status", 2)
}

//发送任务到数据引擎---对接数据引擎3.4.3接口
func (e *Enginejoint) SendTaskToEngine(task orm.Task) (factoryID string, err error) {
	switch task.CurrentStage {
	case 1:
		td, _ := json.Marshal(orm.TaskDescription{
			MsgID:        task.ID,
			Time:         time.Now().Unix(),
			SrcFilePath:  task.VideoPath,
			DestFilePath: task.TransportFilePath,
		})
		sendTask := orm.RealeseTaskToEngine{
			FactoryID: task.TransportFacID,
			Time:      time.Now().Unix(),
			Function:  "生成MP4",
			Value:     string(td),
		}
		msg, _ := json.Marshal(sendTask)
		return sendTask.FactoryID, kafka.SendMessage("factoryFuncInfo", msg)
	case 2:
		td, _ := json.Marshal(orm.TaskDescription{
			MsgID:        task.ID,
			Time:         time.Now().Unix(),
			SrcFilePath:  task.TransportFilePath,
			DestFilePath: task.TrajectoryFilePath,
		})
		sendTask := orm.RealeseTaskToEngine{
			FactoryID: task.TrajectoryFacID,
			Time:      time.Now().Unix(),
			Function:  "生成轨迹文件",
			Value:     string(td),
		}
		msg, _ := json.Marshal(sendTask)
		return sendTask.FactoryID, kafka.SendMessage("factoryFuncInfo", msg)
	case 3:
		td, _ := json.Marshal(orm.TaskDescription02{
			MsgID:                 task.ID,
			Time:                  time.Now().Unix(),
			SrcVideoFilePath:      task.TransportFilePath,
			SrcTrajectoryFilePath: task.TrajectoryFilePath,
			DestImgFilePath:       task.BackgroundFilePath,
		})
		sendTask := orm.RealeseTaskToEngine{
			FactoryID: task.BackgroundFacID,
			Time:      time.Now().Unix(),
			Function:  "生成背景图片",
			Value:     string(td),
		}
		msg, _ := json.Marshal(sendTask)
		return sendTask.FactoryID, kafka.SendMessage("factoryFuncInfo", msg)
	}
	return "", nil
}

//获取处理当前task的factoryID
func (e *Enginejoint) GetCurrentFactoryID(task orm.Task) (factoryID string, err error) {
	switch task.CurrentStage {
	case 1:
		return task.TransportFacID, nil
	case 2:
		return task.TrajectoryFacID, nil
	case 3:
		return task.BackgroundFacID, nil
	case 4:
		return "", errors.New("当前任务已完成")

	}
	return
}

//获取订阅数据的任务id,进度,和错误
func (e *Enginejoint) GetTaskIDAndProgress(data orm.Data01) (up orm.UpdateProgress) {
	perscent := float64(data.CurFrame) / float64(data.AllFrame) * 100
	progress := strconv.FormatFloat(perscent, 'f', 2, 64) + "%"
	up = orm.UpdateProgress{
		TaskID:   data.MsgID,
		Progress: progress,
		DoneFlag: data.CurFrame == data.AllFrame,
		Err:      errors.New(data.ErrCode),
	}
	return
}

//更新任务表:更新current_progress与error
func (e *Enginejoint) Update2(up orm.UpdateProgress) {
	if up.DoneFlag {
		e.DB.Model(&orm.Task{ID: up.TaskID}).Updates(map[string]interface{}{"CurrentStage": gorm.Expr("current_stage + ?", 1), "Status": 1, "CurrentProgress": "waiting new assignment"})
		return
	}
	e.DB.Model(&orm.Task{ID: up.TaskID}).Updates(map[string]interface{}{"ErrorCode": up.Err.Error(), "CurrentProgress": up.Progress})
}
