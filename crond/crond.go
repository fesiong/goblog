package crond

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"kandaoni.com/anqicms/provider"
	"math/rand"
	"time"
)

func Crond() {
	crontab := cron.New(cron.WithSeconds())
	//每天执行
	crontab.AddFunc("@daily", dailyTask)
	// 每天8点执行
	crontab.AddFunc("1 1 8 * * *", daily8HourTask)
	// 每小时执行
	crontab.AddFunc("@hourly", hourlyTask)
	// 每10分钟执行
	crontab.AddFunc("1 */10 * * * *", hourly10MinuteTask)
	// 每分钟执行
	crontab.AddFunc("1 * * * * *", minutelyTask)
	crontab.Start()
}

func dailyTask() {
	//每天执行一次，清理很久的statistic
	cleanStatistics()
	// 每天清理一次回收站内容
	CleanArchives()
	// 每天检查VIP
	CleanUserVip()
	// 每天定期优化表
	optimizeTable()
}

func daily8HourTask() {
	//每天检查一次收录量
	QuerySpiderInclude()
	// 每天8点下发前一天网站数据到邮箱
	SendStatisticsMail()
}

func hourlyTask() {
	// 每小时挖词
	startDigKeywords()
	// 每小时检查一次账号状态
	CheckAuthValid()
}

func hourly10MinuteTask() {
	// 每十分钟检查一次提现
	CollectArticles()
}

func minutelyTask() {
	// 每分钟检查一次需要发布的文章
	PublishPlanContents()
	// 每分钟提现
	CheckWithdrawToWechat()
	// 每分钟定期检查订单
	AutoCheckOrders()
	// 每分钟检查一次时间因子
	UpdateTimeFactor()
	// 每分钟检查一次 AI文章计划
	AiArticlePlan()
}

func optimizeTable() {
	// 需要优化的表
	tables := []string{
		"archives",
		"archive_drafts",
	}

	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		for _, t := range tables {
			w.DB.Exec("OPTIMIZE TABLE `?`", gorm.Expr(t))
		}
	}
}

func startDigKeywords() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.StartDigKeywords(false)
	}
}

func cleanStatistics() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.CleanStatistics()
	}
}

func PublishPlanContents() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.PublishPlanArchives()
	}
}

func CleanArchives() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.CleanArchives()
	}
}

func CollectArticles() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.CollectArticles()
	}
}

func QuerySpiderInclude() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.QuerySpiderInclude()
	}
}

func CheckWithdrawToWechat() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.CheckWithdrawToWechat()
	}
}

func AutoCheckOrders() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.AutoCheckOrders()
	}
}

func CleanUserVip() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.CleanUserVip()
	}
}

func CheckAuthValid() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(600)+1) * time.Second)
	defaultSite := provider.CurrentSite(nil)
	defaultSite.AnqiCheckLogin(false)
}

func UpdateTimeFactor() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.TryToRunTimeFactor()
	}
}

func AiArticlePlan() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.SyncAiArticlePlan()
	}
}

func SendStatisticsMail() {
	websites := provider.GetWebsites()
	for _, w := range websites {
		if !w.Initialed {
			continue
		}
		w.SendStatisticsMail()
	}
}
