package pagination_test

import (
	"testing"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

func TestPagination_Normalize(t *testing.T) {
	tests := []struct {
		name       string
		input      pagination.Pagination
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "zero values receive defaults",
			input:      pagination.Pagination{},
			wantLimit:  pagination.DefaultLimit,
			wantOffset: 0,
		},
		{
			name:       "negative values clamped",
			input:      pagination.Pagination{Limit: -10, Offset: -5},
			wantLimit:  pagination.DefaultLimit,
			wantOffset: 0,
		},
		{
			name:       "exceeds max limit clamped",
			input:      pagination.Pagination{Limit: 500, Offset: 10},
			wantLimit:  pagination.MaxLimit,
			wantOffset: 10,
		},
		{
			name:       "valid custom parameters preserved",
			input:      pagination.Pagination{Limit: 50, Offset: 100},
			wantLimit:  50,
			wantOffset: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Normalize()
			if got.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.wantLimit)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %d, want %d", got.Offset, tt.wantOffset)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	page := pagination.NewPage(150, -5, 42)
	if page.Limit != pagination.MaxLimit {
		t.Errorf("Limit = %d, want %d", page.Limit, pagination.MaxLimit)
	}
	if page.Offset != 0 {
		t.Errorf("Offset = %d, want 0", page.Offset)
	}
	if page.Total != 42 {
		t.Errorf("Total = %d, want 42", page.Total)
	}
}
