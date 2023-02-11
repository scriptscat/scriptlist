package gray_control

import (
	"testing"
	"time"
)

func TestWeight_match(t *testing.T) {
	type fields struct {
		weight    int
		weightDay float64
	}
	type args struct {
		now        time.Time
		n          int
		createtime int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{name: "case1", fields: fields{30, 0}, args: args{time.Now(), 0, 0}, want: true},
		{name: "case2", fields: fields{30, 0}, args: args{time.Now(), 1, 0}, want: true},
		{name: "case3", fields: fields{30, 0}, args: args{time.Now(), 10, 0}, want: true},
		{name: "case4", fields: fields{30, 0}, args: args{time.Now(), 30, 0}, want: true},
		{name: "case5", fields: fields{30, 0}, args: args{time.Now(), 31, 0}, want: false},
		{name: "case6", fields: fields{30, 0}, args: args{time.Now(), 99, 0}, want: false},
		{name: "case7", fields: fields{30, 0}, args: args{time.Now(), 100, 0}, want: true},

		{name: "case1020-69", fields: fields{30, 0}, args: args{time.Unix(0, 0), 69, 1030}, want: false},
		{name: "case1020-70", fields: fields{30, 0}, args: args{time.Unix(0, 0), 70, 1030}, want: true},
		{name: "case1020-99", fields: fields{30, 0}, args: args{time.Unix(0, 0), 99, 1030}, want: true},
		{name: "case1020-1", fields: fields{30, 0}, args: args{time.Unix(0, 0), 1, 1030}, want: false},

		{name: "case-day-1", fields: fields{100, 10}, args: args{time.Unix(86400, 0), 89, 10}, want: false},
		{name: "case-day-2", fields: fields{100, 10}, args: args{time.Unix(86400, 0), 90, 10}, want: true},
		{name: "case-day-3", fields: fields{100, 10}, args: args{time.Unix(86400, 0), 99, 10}, want: true},
		{name: "case-day-4", fields: fields{100, 10}, args: args{time.Unix(86400, 0), 0, 10}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Weight{
				weight:    tt.fields.weight,
				weightDay: tt.fields.weightDay,
			}
			got, err := w.match(tt.args.now, tt.args.n, tt.args.createtime)
			if (err != nil) != tt.wantErr {
				t.Errorf("match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("match() got = %v, want %v", got, tt.want)
			}
		})
	}
}
