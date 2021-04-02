package request

type Pages struct {
	P int `form:"page" binding:"number" label:"页码" json:"page"`
	C int `form:"count" binding:"number" label:"页大小" json:"count"`
}

func NewPage(page, count int) *Pages {
	return &Pages{
		P: page, C: count,
	}
}

func (p *Pages) Page() int {
	if p.P <= 0 {
		return 1
	}
	return p.P
}

func (p *Pages) Size() int {
	if p.C < 10 {
		return 20
	}
	if p.C > 100 {
		return 100
	}
	return p.C
}
