package issue

type ScriptIssue struct {
	ID         int64  `gorm:"column:id" json:"id"`
	ScriptID   int64  `gorm:"column:script_id;type:bigint(20);index:script_id;NOT NULL" json:"script_id"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`
	Title      string `gorm:"column:title;type:varchar(255);NOT NULL" json:"title"`
	Content    string `gorm:"column:content;type:text" json:"content"`
	Labels     string `gorm:"column:labels;type:varchar(255)" json:"labels"`
	Status     int    `gorm:"column:status;type:tinyint(4);default:0;NOT NULL" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)" json:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)" json:"updatetime"`
}

type ScriptIssueComment struct {
	ID         int64  `gorm:"column:id" json:"id"`
	IssueID    int64  `gorm:"column:issue_id;type:bigint(20);index:issue_id;NOT NULL" json:"issue_id"`
	UserID     int64  `gorm:"column:user_id;type:bigint(20);NOT NULL" json:"user_id"`
	Content    string `gorm:"column:content;type:text;NOT NULL" json:"content"`
	Type       int    `gorm:"column:type;type:tinyint(4);default:0;NOT NULL" json:"type"`
	Status     int    `gorm:"column:status;type:tinyint(4);default:0;NOT NULL" json:"status"`
	Createtime int64  `gorm:"column:createtime;type:bigint(20)" json:"createtime"`
	Updatetime int64  `gorm:"column:updatetime;type:bigint(20)" json:"updatetime"`
}

type IssueLabel struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
