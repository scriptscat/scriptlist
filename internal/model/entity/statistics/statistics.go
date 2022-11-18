package statistics

type StatisticsDownload struct {
	ID              int64  `gorm:"column:id" json:"id" form:"id"`
	UserId          int64  `gorm:"column:user_id" json:"user_id" form:"user_id"`
	ScriptId        int64  `gorm:"column:script_id;index:script_id;index:script_time" json:"script_id" form:"script_id"`
	ScriptCodeId    int64  `gorm:"column:script_code_id;index:script_code_id;index:script_code_time" json:"script_code_id" form:"script_code_id"`
	Ip              string `gorm:"column:ip;type:varchar(100)" json:"ip" form:"ip"`
	Ua              string `gorm:"column:ua;type:text" json:"ua" form:"ua"`
	StatisticsToken string `gorm:"column:statistics_token;type:text" json:"statistics_token" form:"statistics_token"`
	Createtime      int64  `gorm:"column:createtime;index:script_time;index:script_code_time" json:"createtime" form:"createtime"`
}

type StatisticsUpdate StatisticsDownload

type StatisticsPageView StatisticsDownload

type Statistics interface {
	GetUserId() int64
	GetScriptId() int64
	GetIp() string
	GetUa() string
	GetStatisticsToken() string
}

func (s *StatisticsDownload) GetUserId() int64 {
	return s.UserId
}

func (s *StatisticsDownload) GetScriptId() int64 {
	return s.ScriptId
}

func (s *StatisticsDownload) GetIp() string {
	return s.Ip
}

func (s *StatisticsDownload) GetUa() string {
	return s.Ua
}

func (s *StatisticsDownload) GetStatisticsToken() string {
	return s.StatisticsToken
}

func (s *StatisticsUpdate) GetUserId() int64 {
	return s.UserId
}

func (s *StatisticsUpdate) GetScriptId() int64 {
	return s.ScriptId
}

func (s *StatisticsUpdate) GetIp() string {
	return s.Ip
}

func (s *StatisticsUpdate) GetUa() string {
	return s.Ua
}

func (s *StatisticsUpdate) GetStatisticsToken() string {
	return s.StatisticsToken
}

func (s *StatisticsPageView) GetUserId() int64 {
	return s.UserId
}

func (s *StatisticsPageView) GetScriptId() int64 {
	return s.ScriptId
}

func (s *StatisticsPageView) GetIp() string {
	return s.Ip
}

func (s *StatisticsPageView) GetUa() string {
	return s.Ua
}

func (s *StatisticsPageView) GetStatisticsToken() string {
	return s.StatisticsToken
}
