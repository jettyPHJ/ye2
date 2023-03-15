package task

import (
	"encoding/json"
	ch "file_handling/moudle/chaos_corner"
	e "file_handling/moudle/engine_joint"
	t "file_handling/moudle/task_release"
	"file_handling/orm"
	r "file_handling/utils/gin"
	k "file_handling/utils/kafka"
	"file_handling/utils/mysql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	v := r.Router.Group("/task")
	{
		v.POST("", post)
	}
	go enginejoint1()
	go enginejoint2()
	go enginejoint3()
	go supervise()
}

func post(c *gin.Context) {
	videoMeta := orm.Video{}
	ch.RespondWithError(c, c.ShouldBindJSON(&videoMeta))
	inserter := &t.DBInserter{}
	ch.RespondWithError(c, inserter.InsertDefault(videoMeta))
	c.JSON(200, gin.H{
		"msg:": "操作成功",
	})
}

//查询任务条件为status：1（waiting）
func enginejoint1() {
	joint := &e.Enginejoint{DB: mysql.DataBase}
	for {
		tasks := joint.SearchForTask(1)
		for _, v := range tasks {
			task := v
			go func(orm.Task, *e.Enginejoint) {
				factoryInfoes := joint.SearchForFactoryID(task)
				factory, err := joint.SelectOneFactory(factoryInfoes)
				if err != nil {
					return
				}
				if joint.CheckFactoryFree(*factory, task) == "" {
					joint.Update1(task, *factory)
				}
			}(task, joint)
		}
		time.Sleep(3 * time.Second)
	}
}

//查询任务条件为status：2（ready）,向数据引擎发布task
func enginejoint2() {
	joint := &e.Enginejoint{DB: mysql.DataBase}
	for {
		tasks := joint.SearchForTask(2)
		for _, task := range tasks {
			if facID, err := joint.SendTaskToEngine(task); err != nil {
				panic(err)
			} else {
				time.Sleep(1 * time.Second)
				//查询任务是否已被接收，若被接收，更新status：3，若没有，更新status:1
				if joint.CheckFactoryFree(orm.FactoryInfo{FactoryID: facID}, task) != task.ID {
					mysql.DataBase.Model(&task).Update("status", 1)
					continue
				}
				mysql.DataBase.Model(&task).Update("status", 3)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

//订阅实时数据，更新数据库的任务进度
func enginejoint3() {
	joint := &e.Enginejoint{DB: mysql.DataBase}
	msgQueue := make(chan []byte, 100)
	go k.ConsumeMsg("factoryDataUp", msgQueue)
	for {
		msgBytes := <-msgQueue
		sd := orm.SubscriptionData{}
		if err := json.Unmarshal(msgBytes, &sd); err != nil {
			continue
		}
		data := orm.Data01{}
		if err := json.Unmarshal([]byte(sd.Data), &data); err != nil {
			continue
		}
		fmt.Println(data)
		up := joint.GetTaskIDAndProgress(data)
		joint.Update2(up)
	}
}

//监控任务进行中的异常，如果某任务处于status:3 一定时间没有更新进度，则放弃（重新将status置为1）
func supervise() {
	joint := &e.Enginejoint{DB: mysql.DataBase}
	recordMap1 := map[string]struct{}{} //用于记录上一次status:3的task数据id
	for {
		tasks := joint.SearchForTask(3)
		recordMap2 := map[string]struct{}{} //用于记录当前status:3的task数据id
		for _, task := range tasks {
			id := task.ID + fmt.Sprint(task.CurrentStage) + task.CurrentProgress
			recordMap2[id] = struct{}{}
			_, ok := recordMap1[id]
			if ok { //如果上一次就记录了该id,那么这个task的status改为1，重来一遍
				mysql.DataBase.Model(&task).Update("status", 1)
			}
		}
		recordMap1 = recordMap2 //更新recordMap1
		time.Sleep(30 * time.Second)
	}
}
