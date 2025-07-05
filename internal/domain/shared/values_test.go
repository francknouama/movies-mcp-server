package shared

import (
	"testing"
)

func TestNewMovieID(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid positive ID",
			value:   1,
			wantErr: false,
		},
		{
			name:    "valid large ID",
			value:   999999,
			wantErr: false,
		},
		{
			name:    "valid zero ID",
			value:   0,
			wantErr: false,
		},
		{
			name:    "invalid negative ID",
			value:   -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewMovieID(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMovieID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id.Value() != tt.value {
				t.Errorf("NewMovieID() value = %v, want %v", id.Value(), tt.value)
			}
		})
	}
}

func TestMovieID_IsZero(t *testing.T) {
	zeroID := MovieID{}
	if !zeroID.IsZero() {
		t.Error("Expected zero MovieID to return true for IsZero()")
	}

	validID, _ := NewMovieID(1)
	if validID.IsZero() {
		t.Error("Expected valid MovieID to return false for IsZero()")
	}
}

func TestNewActorID(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid positive ID",
			value:   1,
			wantErr: false,
		},
		{
			name:    "valid zero ID",
			value:   0,
			wantErr: false,
		},
		{
			name:    "invalid negative ID",
			value:   -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewActorID(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActorID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id.Value() != tt.value {
				t.Errorf("NewActorID() value = %v, want %v", id.Value(), tt.value)
			}
		})
	}
}

func TestNewRating(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{
			name:    "valid minimum rating",
			value:   0.0,
			wantErr: false,
		},
		{
			name:    "valid maximum rating",
			value:   10.0,
			wantErr: false,
		},
		{
			name:    "valid decimal rating",
			value:   7.5,
			wantErr: false,
		},
		{
			name:    "invalid negative rating",
			value:   -1.0,
			wantErr: true,
		},
		{
			name:    "invalid too high rating",
			value:   10.1,
			wantErr: true,
		},
		{
			name:    "invalid much too high rating",
			value:   100.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating, err := NewRating(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRating() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && rating.Value() != tt.value {
				t.Errorf("NewRating() value = %v, want %v", rating.Value(), tt.value)
			}
		})
	}
}

func TestRating_IsZero(t *testing.T) {
	zeroRating := Rating{}
	if !zeroRating.IsZero() {
		t.Error("Expected zero Rating to return true for IsZero()")
	}

	validRating, _ := NewRating(7.5)
	if validRating.IsZero() {
		t.Error("Expected valid Rating to return false for IsZero()")
	}
}

func TestNewYear(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid year - first movie",
			value:   1888,
			wantErr: false,
		},
		{
			name:    "valid year - current",
			value:   2024,
			wantErr: false,
		},
		{
			name:    "valid year - near future",
			value:   2030,
			wantErr: false,
		},
		{
			name:    "invalid year - too old",
			value:   1887,
			wantErr: true,
		},
		{
			name:    "invalid year - too far future",
			value:   2050,
			wantErr: true,
		},
		{
			name:    "invalid year - negative",
			value:   -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, err := NewYear(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewYear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && year.Value() != tt.value {
				t.Errorf("NewYear() value = %v, want %v", year.Value(), tt.value)
			}
		})
	}
}
