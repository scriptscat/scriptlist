package script_entity

type ScriptDomain struct {
	ID            int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	Domain        string `gorm:"column:domain;type:varchar(255);index:domain_script,unique"`
	DomainReverse string `gorm:"column:domain_reverse;type:varchar(255);index:domain_reverse"` // 域名反转, 方便查询
	ScriptID      int64  `gorm:"column:script_id;type:bigint(20);index:domain_script,unique;index:script_id"`
	ScriptCodeID  int64  `gorm:"column:script_code_id;type:bigint(20)"`
	Createtime    int64  `gorm:"column:createtime;type:bigint(20)"`
}
