package task

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"platform-backend/config"
	db "platform-backend/db"
	"platform-backend/dto/enum"
	"platform-backend/models"
	"platform-backend/repository"
	"platform-backend/utils/log"
)

type Task interface {
	Exec()
	GetPolicy() *models.Policy
}

type taskDetail struct {
	entryID cron.EntryID
	policy  *models.Policy
}

var taskEntryMap = make(map[int64]*taskDetail)

func RunTask() error {

	cronScheduler := cron.New(cron.WithChain(cron.SkipIfStillRunning(log.SystemLogger)), cron.WithLocation(time.UTC))

	gapTime := config.MyConfig.CronGapTime
	if config.MyConfig.CronGapTime < 5 {
		gapTime = 5
	}

	ticker := time.NewTicker(time.Duration(gapTime) * time.Second)
	defer ticker.Stop()

	t := TimeOnlineTask{}
	taskID, err := cronScheduler.AddFunc("0/1 * * * *", t.Exec)
	if err != nil {
		return err
	}
	defer cronScheduler.Remove(taskID)

	err = scheduleCronTasks(cronScheduler)
	if err != nil {
		return err
	}

	for range ticker.C {
		if err := scheduleCronTasks(cronScheduler); err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "重新加载任务失败:%s\n", err.Error())
		} else {
			_, _ = fmt.Fprintf(log.SystemLogger, "加载任务成功\n")
		}
		cronScheduler.Start()
	}

	return nil
}

func fetchCronTasks() ([]Task, error) {
	policyRepo := repository.NewPolicyRepository(db.MysqlDB.DB())
	policies, err := policyRepo.FindAll(context.Background(), map[string]interface{}{
		"status": enum.PolicyStatusEnable,
	}, false)

	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	for _, policy := range policies {
		switch policy.ActionType {
		case enum.ActionTypeOn:
			t := PowerOnTask{BaseTaskImpl{Policy: &policy}}
			tasks = append(tasks, &t)
		case enum.ActionTypeOff:
			t := PowerOffTask{BaseTaskImpl{Policy: &policy}}
			tasks = append(tasks, &t)
		default:
			continue
		}
	}

	return tasks, nil
}

// scheduleCronTasks 根据数据库中的 cron 表达式安排任务
func scheduleCronTasks(cronScheduler *cron.Cron) error {
	tasks, err := fetchCronTasks()
	if err != nil {
		return err
	}

	newTaskMap := make(map[int64]struct{})
	for _, task := range tasks {
		policy := task.GetPolicy()
		newTaskMap[policy.ID] = struct{}{}

		if v, ok := taskEntryMap[policy.ID]; ok {
			if v.policy.UpdatedAt != policy.UpdatedAt {
				// 任务变更了
				_, _ = fmt.Fprintf(log.SystemLogger, "重新加载 %s 策略\n", policy.PolicyName)
				cronScheduler.Remove(v.entryID)
			} else {
				//没有变化
				continue
			}
		}
		executeTime, err := time.Parse("15:04:05", policy.ExecuteTime)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "解析 %s 策略时间失败, %s\n", policy.PolicyName, err.Error())
			continue
		}
		//
		var CronTime string
		if policy.ExecuteType == enum.TypeOnce {
			CronTime = fmt.Sprintf("%d %d %d %d ?",
				executeTime.Minute(),
				executeTime.Hour(),
				policy.StartDate.Day(),
				policy.StartDate.Month())
		} else if policy.ExecuteType == enum.TypeEveryDay {
			CronTime = fmt.Sprintf("%d %d * ? ? ",
				executeTime.Minute(),
				executeTime.Hour())
		} else if policy.ExecuteType == enum.TypeDayOfWeek {
			weekstr := ""
			havestr := false
			if policy.DayOfWeek&0x02 == 0x02 {
				weekstr = fmt.Sprintf("%s%s", weekstr, "1")
				havestr = true
			}
			if policy.DayOfWeek&0x04 == 0x04 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "2")
				havestr = true
			}
			if policy.DayOfWeek&0x08 == 0x08 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "3")
				havestr = true
			}
			if policy.DayOfWeek&0x10 == 0x10 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "4")
				havestr = true
			}
			if policy.DayOfWeek&0x20 == 0x20 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "5")
				havestr = true
			}
			if policy.DayOfWeek&0x40 == 0x40 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "6")
				havestr = true
			}
			if policy.DayOfWeek&0x80 == 0x80 {
				if havestr {
					weekstr = fmt.Sprintf("%s,", weekstr)
				}
				weekstr = fmt.Sprintf("%s%s", weekstr, "7")
				havestr = true
			}
			CronTime = fmt.Sprintf("%d %d ? ? %s",
				executeTime.Minute(),
				executeTime.Hour(),
				weekstr)
		}

		entryID, err := cronScheduler.AddFunc(CronTime, task.Exec)
		if err != nil {
			_, _ = fmt.Fprintf(log.SystemLogger, "任务调度失败: %s %v\n", err.Error(), task)
			continue
		}

		taskEntryMap[policy.ID] = &taskDetail{
			entryID: entryID,
			policy:  policy,
		}

	}

	// 去除已删除的
	for key, v := range taskEntryMap {
		if _, ok := newTaskMap[key]; ok {
			continue
		}

		delete(taskEntryMap, key)
		cronScheduler.Remove(v.entryID)
	}

	return nil
}
