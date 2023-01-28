package resource_entity

type Resource struct {
	ID          int64  `gorm:"column:id;type:bigint(20);not null;primary_key"`
	ResourceID  string `gorm:"column:resource_id;type:varchar(255);not null;index:resource,unique"`
	UserID      int64  `gorm:"column:user_id;type:bigint(20);not null"`
	LinkID      int64  `gorm:"column:link_id;type:bigint(20)"`
	Comment     string `gorm:"column:comment;type:varchar(255)"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	Path        string `gorm:"column:path;type:varchar(255);not null"`
	ContentType string `gorm:"column:content_type;type:varchar(255)"`
	Status      int    `gorm:"column:status;type:int(11);not null"`
	Createtime  int64  `gorm:"column:createtime;type:bigint(20)"`
}
