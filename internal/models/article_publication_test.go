package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestArticleResponseFirstPublicationTime(t *testing.T) {
	published := time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC)
	for _, tt := range []struct {
		name string
		at   *time.Time
		want string
	}{
		{name: "never public", want: "null"},
		{name: "first public time", at: &published, want: `"2020-02-03T04:05:06Z"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal((&Article{PublishAt: tt.at}).ToResponse())
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(body, &fields); err != nil {
				t.Fatal(err)
			}
			if string(fields["publish_at"]) != tt.want {
				t.Fatalf("publish_at=%s, want %s", fields["publish_at"], tt.want)
			}
		})
	}
}
