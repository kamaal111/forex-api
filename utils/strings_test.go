package utils

import "testing"

func TestArrayContains(t *testing.T) {
	tests := []struct {
		name   string
		array  []string
		search string
		want   bool
	}{
		{
			name:   "element exists in array",
			array:  []string{"EUR", "USD", "GBP"},
			search: "USD",
			want:   true,
		},
		{
			name:   "element does not exist in array",
			array:  []string{"EUR", "USD", "GBP"},
			search: "JPY",
			want:   false,
		},
		{
			name:   "empty array",
			array:  []string{},
			search: "EUR",
			want:   false,
		},
		{
			name:   "element is first in array",
			array:  []string{"EUR", "USD", "GBP"},
			search: "EUR",
			want:   true,
		},
		{
			name:   "element is last in array",
			array:  []string{"EUR", "USD", "GBP"},
			search: "GBP",
			want:   true,
		},
		{
			name:   "case sensitive match",
			array:  []string{"EUR", "USD", "GBP"},
			search: "eur",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayContains(tt.array, tt.search)
			if got != tt.want {
				t.Errorf("ArrayContains(%v, %q) = %v, want %v", tt.array, tt.search, got, tt.want)
			}
		})
	}
}

func TestArrayContainsInt(t *testing.T) {
	tests := []struct {
		name   string
		array  []int
		search int
		want   bool
	}{
		{
			name:   "integer exists in array",
			array:  []int{1, 2, 3, 4, 5},
			search: 3,
			want:   true,
		},
		{
			name:   "integer does not exist in array",
			array:  []int{1, 2, 3, 4, 5},
			search: 6,
			want:   false,
		},
		{
			name:   "empty array",
			array:  []int{},
			search: 1,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayContains(tt.array, tt.search)
			if got != tt.want {
				t.Errorf("ArrayContains(%v, %d) = %v, want %v", tt.array, tt.search, got, tt.want)
			}
		})
	}
}
