package statistics_entity

// StatisticsCollect 统计收集
type StatisticsCollect struct {
	SessionID     string `gorm:"column:session_id;type:varchar(128);not null;primary_key;comment:会话id"`
	ScriptID      int64  `gorm:"column:script_id;type:bigint(20);not null;comment:脚本id"`
	VisitorID     string `gorm:"column:visitor_id;type:varchar(128);not null;comment:访客id"`
	OperationHost string `gorm:"column:operation_host;type:string;not null;comment:操作域名"`
	OperationPage string `gorm:"column:operation_page;type:string;not null;comment:操作页面"`
	Duration      int32  `gorm:"column:duration;type:int(10);not null;comment:停留时长"`
	VisitTime     int64  `gorm:"column:visit_time;type:bigint(20);not null;comment:访问时间"`
	ExitTime      int64  `gorm:"column:exit_time;type:bigint(20);not null;comment:退出时间"`
}

// StatisticsVisitor 访客统计
type StatisticsVisitor struct {
	ScriptID    int64  `gorm:"column:script_id;type:bigint(20);not null;comment:脚本id"`
	VisitorID   string `gorm:"column:visitor_id;type:varchar(128);not null;comment:访客id"`
	UA          string `gorm:"column:ua;type:string;not null"`
	IP          string `gorm:"column:ip;type:varchar(128);not null"`
	Version     string `gorm:"column:version;type:varchar(128);not null;comment:版本"`
	InstallPage string `gorm:"column:install_page;type:string;not null;comment:安装页面"`
	InstallHost string `gorm:"column:install_host;type:string;not null;comment:安装域名"`
}
