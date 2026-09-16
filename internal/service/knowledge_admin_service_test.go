package service

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestAdminPrivateDocSerialization(t *testing.T) {
	for _, v := range []interface{}{AdminDocMetadata{Title: "private"}, AdminDocDetail{AdminDocMetadata: AdminDocMetadata{Title: "draft", Status: 0}}} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), `"content"`) || strings.Contains(string(b), `"content_html"`) {
			t.Fatalf("body field leaked: %s", b)
		}
	}
}
func TestKnowledgeTrendZeroFill(t *testing.T) {
	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)
	rows := fillKnowledgeTrend(now, []KnowledgeTrend{{"2026-09-10", 4}, {"2026-09-15", 2}})
	if len(rows) != 7 || rows[0].Day != "2026-09-09" || rows[6].Day != "2026-09-15" {
		t.Fatalf("invalid dates: %+v", rows)
	}
	for i, r := range rows {
		want := int64(0)
		if i == 1 {
			want = 4
		}
		if i == 6 {
			want = 2
		}
		if r.Count != want {
			t.Errorf("day %s count=%d want %d", r.Day, r.Count, want)
		}
	}
}
func TestKnowledgeQuotaRejectsInvalidKeysBeforeDatabase(t *testing.T) {
	for _, v := range []map[string]int{{"site_name": 1}, {QuotaMaxDocs: -1}, {QuotaMaxStorage: -5}} {
		if NewKnowledgeAdminService().SaveQuota(v) == nil {
			t.Fatalf("accepted invalid quota: %v", v)
		}
	}
}
