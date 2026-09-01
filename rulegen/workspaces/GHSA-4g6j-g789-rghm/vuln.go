package main

	"fmt"
	"slices"
	"strings"

	"github.com/jinzhu/copier"

	pb "github.com/nezhahq/nezha/proto"
)

type CronClass struct {
	class[uint64, *model.Cron]
	*cron.Cron
}

func NewCronClass() *CronClass {
			list:       list,
			sortedList: sortedList,
		},
		Cron: cronx,
	}
}

	delete(c.list, cr.ID)
	c.list[cr.ID] = cr
	c.listMu.Unlock()

	c.sortList()
}
		delete(c.list, id)
	}
	c.listMu.Unlock()

	c.sortList()
}
	return cr.UserID == triggerOwner || userIsAdmin(triggerOwner)
}

func ManualTrigger(cr *model.Cron) {
	CronTrigger(cr)()
}
					return
				}
				if s.TaskStream != nil {
					s.TaskStream.Send(&pb.Task{
						Id:   cr.ID,
						Data: cr.Command,
						Type: model.TaskTypeCommand,
					})
				} else {
					// 保存当前服务器状态信息
					curServer := model.Server{}
	return true
}

// worker 服务监控的实际工作流程
func (ss *ServiceSentinel) worker() {
	// 从服务状态汇报管道获取汇报的服务数据
	for r := range ss.serviceReportChannel {
		css, _ := ss.Get(r.Data.GetId())
		if css == nil || css.ID == 0 {
			log.Printf("NEZHA>> Incorrect service monitor report %+v", r)
			continue
		}
		css = nil

		mh := r.Data
		if mh.Type == model.TaskTypeTCPPing || mh.Type == model.TaskTypeICMPPing {
			ss.serviceCurrentStatusData[mh.GetId()].result = ss.serviceCurrentStatusData[mh.GetId()].result[:0]
		}

		cs, _ := ss.Get(mh.GetId())
		m := ServerShared.GetList()
		// 延迟报警
		if mh.Delay > 0 {
