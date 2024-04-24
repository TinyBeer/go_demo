package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestAddUpper(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{},
		{
			name: "累加到10",
			args: args{
				n: 10,
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddUpper(tt.args.n); got != tt.want {
				t.Errorf("AddUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		str string
		sep string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "单字符分隔",
			args: args{
				str: "a:b:c:d:",
				sep: ":",
			},
			want: []string{"a", "b", "c", "d", ""},
		},
		{
			name: "多字符分隔",
			args: args{
				str: "a::b::c::d::",
				sep: "::",
			},
			want: []string{"a", "b", "c", "d", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.args.str, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSplit(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Split("a:b:c:d", ":")
	}
}

func BenchmarkStringsSplit(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		strings.Split("a:b:c:d", ":")
	}
}
