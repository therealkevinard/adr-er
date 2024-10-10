package utils

import (
	"testing"
)

func TestPadValue(t *testing.T) {
	type args struct {
		value int
		width int
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "",
			want: "001",
			args: args{width: 3, value: 1},
		},
		{
			name: "",
			want: "011",
			args: args{width: 3, value: 11},
		},
		{
			name: "",
			want: "111",
			args: args{width: 3, value: 111},
		},
		{
			name: "",
			want: "1111",
			args: args{width: 3, value: 1111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadValue(tt.args.value, tt.args.width)
			if got != tt.want {
				t.Errorf("PadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		want string
		args string
	}{
		{
			name: "basic",
			want: "fizzy-pop",
			args: "Fizzy Pop",
		},
		{
			name: "many spaces",
			want: "fizzy-pop",
			args: "Fizzy                 Pop",
		},
		{
			name: "leading-trailing whsp",
			want: "fizzy-pop",
			args: "   Fizzy Pop   ",
		},
		{
			name: "unsafe characters",
			want: "fizzy-pop",
			args: "Fiz*zy P\\o/p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.args)
			if got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugifySlice(t *testing.T) {
	tests := []struct {
		name string
		want string
		args []string
	}{
		{
			name: "basic",
			want: "fizzy-pop",
			args: []string{"fizzy", "pop"},
		},
		{
			name: "with padded num",
			want: "001-fizzy-pop",
			args: []string{"001", "fizzy", "pop"},
		},
		{
			name: "pre-slugged",
			want: "001-fizzy-pop",
			args: []string{"001", "fizzy-pop"},
		},
		{
			name: "with-zero-value",
			want: "001-fizzy-pop",
			args: []string{"", "", "", "", "001", "fizzy", "pop", "", "", "", ""},
		},
		{
			name: "with-bad chars",
			want: "001-fizzy-pop",
			args: []string{"-", "-", "//\\", "*", "001", "fizzy", "pop", "(", "#", "&", "%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlugifySlice(tt.args...); got != tt.want {
				t.Errorf("SlugifySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
