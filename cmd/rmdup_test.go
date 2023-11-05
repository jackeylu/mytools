package cmd

import (
	"testing"
	"time"
)

func TestRmdup1(t *testing.T) {
	testCases := []struct {
		desc   string
		emails []EmailInfo
		want   []EmailInfo
	}{
		{
			desc: "1 and 3 duplicated",
			emails: []EmailInfo{
				// Create sample email info
				{
					SeqNum:  1,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "sender1@example.com",
					To:      []string{"receiver1@example.com"},
					Subject: "This is Project",
					Attachments: []string{
						"attachment1.txt",
						"attachment2.txt",
					},
				},
				{
					SeqNum:  1,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "sender2@example.com",
					To:      []string{"receiver1@example.com"},
					Subject: "This is Project",
					Attachments: []string{
						"attachment1.txt",
						"attachment2.txt",
					},
				},
				{
					SeqNum:  199,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "sender1@example.com",
					To:      []string{"receiver1@example.com"},
					Subject: "This is Project",
					Attachments: []string{
						"attachment1.txt",
						"attachment2.txt",
					},
				},
			},
			want: []EmailInfo{
				// Create sample email info
				{
					SeqNum:  1,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "sender1@example.com",
					To:      []string{"receiver1@example.com"},
					Subject: "This is Project",
					Attachments: []string{
						"attachment1.txt",
						"attachment2.txt",
					},
				},
				{
					SeqNum:  199,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "sender2@example.com",
					To:      []string{"receiver1@example.com"},
					Subject: "This is Project",
					Attachments: []string{
						"attachment1.txt",
						"attachment2.txt",
					},
				},
			},
		},
		{
			desc: "2 and 3 duplicated",
			emails: []EmailInfo{
				// Create sample email info
				{
					SeqNum:  500,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "488631603@qq.com",
					To:      []string{"daiyongfen@anncare.cn", "378554759@qq.com"},
					Subject: "应收理赔款批次20230814/数量1/总金额468.120元/耗时107毫秒",
				},
				{
					SeqNum:  501,
					Date:    time.Date(2019, 1, 1, 0, 10, 0, 0, time.UTC),
					From:    "488631603@qq.com",
					To:      []string{"daiyongfen@anncare.cn", "378554759@qq.com"},
					Subject: "应收服务费批次20230814/数量2/总金额1036.000元/耗时163毫秒",
				},
			},
			want: []EmailInfo{
				// Create sample email info
				{
					SeqNum:  500,
					Date:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					From:    "488631603@qq.com",
					To:      []string{"daiyongfen@anncare.cn", "378554759@qq.com"},
					Subject: "应收理赔款批次20230814/数量1/总金额468.120元/耗时107毫秒",
				},
				{
					SeqNum:  501,
					Date:    time.Date(2019, 1, 1, 0, 10, 0, 0, time.UTC),
					From:    "488631603@qq.com",
					To:      []string{"daiyongfen@anncare.cn", "378554759@qq.com"},
					Subject: "应收服务费批次20230814/数量2/总金额1036.000元/耗时163毫秒",
				},
			},
		},
	}

	equal := func(a, b []EmailInfo) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if !a[i].Equals(b[i]) {
				return false
			}
		}
		return true
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := rmdup(tC.emails); !equal(got, tC.want) {
				t.Errorf("rmdup() = %v, want %v", got, tC.want)
			}
		})
	}
}
