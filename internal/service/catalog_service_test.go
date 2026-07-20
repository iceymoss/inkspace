package service

import (
	"reflect"
	"testing"

	"github.com/iceymoss/inkspace/internal/models"
)

func uintPtr(value uint) *uint { return &value }

func TestCatalogDescendantIDs(t *testing.T) {
	catalogs := []models.Catalog{
		{ID: 1},
		{ID: 2, ParentID: uintPtr(1)},
		{ID: 3, ParentID: uintPtr(2)},
		{ID: 4, ParentID: uintPtr(1)},
		{ID: 5},
	}

	got := catalogDescendantIDs(1, catalogs)
	want := []uint{1, 4, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("catalogDescendantIDs() = %v, want %v", got, want)
	}
}

func TestWouldCreateCatalogCycle(t *testing.T) {
	catalogs := []models.Catalog{
		{ID: 1},
		{ID: 2, ParentID: uintPtr(1)},
		{ID: 3, ParentID: uintPtr(2)},
		{ID: 4},
	}

	tests := []struct {
		name     string
		parentID *uint
		want     bool
	}{
		{name: "root", parentID: nil, want: false},
		{name: "self", parentID: uintPtr(1), want: true},
		{name: "child", parentID: uintPtr(2), want: true},
		{name: "deep descendant", parentID: uintPtr(3), want: true},
		{name: "unrelated", parentID: uintPtr(4), want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := wouldCreateCatalogCycle(1, test.parentID, catalogs); got != test.want {
				t.Fatalf("wouldCreateCatalogCycle() = %v, want %v", got, test.want)
			}
		})
	}
}
