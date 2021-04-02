package utils

type Pages struct {
	page  int64 `form:"page" binding:"number" label:"页码" json:"page"`
	count int64 `form:"count" binding:"number" label:"页大小" json:"count"`
}

func (p *Pages) Page() int64 {
	if p.page <= 0 {
		return 1
	}
	return p.page
}

func (p *Pages) Size() int64 {
	if p.count < 10 {
		return 20
	}
	if p.count > 100 {
		return 100
	}
	return p.count
}
