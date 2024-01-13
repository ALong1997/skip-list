package skip_list

import (
"golang.org/x/exp/constraints"
"reflect"
"testing"
)

func TestKvPair_Key(t *testing.T) {
	type testCase[O constraints.Ordered, T any] struct {
		name    string
		kv      KvPair[O, T]
		wantKey O
	}
	tests := []testCase[int, int]{
		{
			name:    "TestKvPair_Key 1",
			kv:      KvPair[int, int]{1, 1},
			wantKey: 1,
		},
		{
			name:    "TestKvPair_Key 2",
			kv:      KvPair[int, int]{2, 1},
			wantKey: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.kv.Key(); !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Key() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}

	tests2 := []testCase[string, int]{
		{
			name:    "TestKvPair_Key 3",
			kv:      KvPair[string, int]{"int", 1},
			wantKey: "int",
		},
		{
			name:    "TestKvPair_Key 4",
			kv:      KvPair[string, int]{"string", 1},
			wantKey: "string",
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.kv.Key(); !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Key() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}

	tests3 := []testCase[float64, int]{
		{
			name:    "TestKvPair_Key 5",
			kv:      KvPair[float64, int]{0.01, 1},
			wantKey: 0.01,
		},
		{
			name:    "TestKvPair_Key 6",
			kv:      KvPair[float64, int]{88, 1},
			wantKey: 88,
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.kv.Key(); !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Key() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestKvPair_Val(t *testing.T) {
	type testCase[O constraints.Ordered, T any] struct {
		name    string
		kv      KvPair[O, T]
		wantVal T
	}
	tests := []testCase[int, int]{
		{
			name:    "TestKvPair_Val 1",
			kv:      KvPair[int, int]{1, 1},
			wantVal: 1,
		},
		{
			name:    "TestKvPair_Val 2",
			kv:      KvPair[int, int]{2, 1},
			wantVal: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVal := tt.kv.Val(); !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Val() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}

	tests2 := []testCase[string, int]{
		{
			name:    "TestKvPair_Val 3",
			kv:      KvPair[string, int]{"int", 1},
			wantVal: 1,
		},
		{
			name:    "TestKvPair_Val 4",
			kv:      KvPair[string, int]{"string", 1},
			wantVal: 1,
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.kv.Val(); !reflect.DeepEqual(gotKey, tt.wantVal) {
				t.Errorf("Val() = %v, want %v", gotKey, tt.wantVal)
			}
		})
	}

	tests3 := []testCase[float64, int]{
		{
			name:    "TestKvPair_Val 5",
			kv:      KvPair[float64, int]{0.01, 1},
			wantVal: 1,
		},
		{
			name:    "TestKvPair_Val 6",
			kv:      KvPair[float64, int]{88, 1},
			wantVal: 1,
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.kv.Val(); !reflect.DeepEqual(gotKey, tt.wantVal) {
				t.Errorf("Val() = %v, want %v", gotKey, tt.wantVal)
			}
		})
	}
}

func Test_newKvPair(t *testing.T) {
	type args[O constraints.Ordered, T any] struct {
		key O
		val T
	}
	type testCase[O constraints.Ordered, T any] struct {
		name string
		args args[O, T]
		want *KvPair[O, T]
	}
	tests := []testCase[int, int]{
		{
			name: "Test_newKvPair 1",
			args: args[int, int]{1, 1},
			want: &KvPair[int, int]{1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKvPair(tt.args.key, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKvPair() = %v, want %v", got, tt.want)
			}
		})
	}

	tests2 := []testCase[string, int]{
		{
			name: "Test_newKvPair 2",
			args: args[string, int]{"int", 1},
			want: &KvPair[string, int]{"int", 1},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKvPair(tt.args.key, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKvPair() = %v, want %v", got, tt.want)
			}
		})
	}

	tests3 := []testCase[float64, int]{
		{
			name: "Test_newKvPair 2",
			args: args[float64, int]{0.01, 1},
			want: &KvPair[float64, int]{0.01, 1},
		},
	}
	for _, tt := range tests3 {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKvPair(tt.args.key, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKvPair() = %v, want %v", got, tt.want)
			}
		})
	}
}
