package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDocResponseUsesLegacyDefaults(t *testing.T) {
	response := (&Doc{}).ToResponse()
	if response.Kind != DocKindMarkdown || response.Revision != 1 {
		t.Fatalf("legacy defaults = kind %q revision %d, want markdown revision 1", response.Kind, response.Revision)
	}

	version := (&DocVersion{}).ToResponse()
	if version.Kind != DocKindMarkdown || version.Revision != 1 {
		t.Fatalf("version legacy defaults = kind %q revision %d, want markdown revision 1", version.Kind, version.Revision)
	}
}

func TestDocDetailResponseDoesNotExposeStorageFields(t *testing.T) {
	body, err := json.Marshal(&DocDetailResponse{Attachment: &DocAttachmentDetail{ID: 3, FileName: "safe.pdf"}})
	if err != nil {
		t.Fatalf("marshal detail response: %v", err)
	}
	jsonBody := string(body)
	for _, forbidden := range []string{"attachment_id", "published_attachment_id", "file_path", "url", "owner_id"} {
		if strings.Contains(jsonBody, forbidden) {
			t.Fatalf("detail response exposed %q: %s", forbidden, jsonBody)
		}
	}
}
