package reports

import "time"

type Getter interface {
	GetReport() Report
}
type Implementation struct{}

func NewGetter() *Implementation {
	return &Implementation{}
}

func (i *Implementation) GetReport() Report {
	return Report{
		Data: "report from " + time.Now().String(),
	}
}
