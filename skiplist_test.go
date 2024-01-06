package skiplist

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

func TestSkipList_Level(t *testing.T) {
	type testCase[O constraints.Ordered, T any] struct {
		name string
		sl   *SkipList[O, T]
		want int32
	}
	var sl *SkipList[int, int]
	tt := testCase[int, int]{
		name: "TestSkipList_Level 1",
		sl:   sl,
	}
	t.Run(tt.name, func(t *testing.T) {
		t.Logf("Level() = %v", tt.sl.Level())
	})

	sl = NewSkipList[int, int](10, false)
	t.Run(tt.name, func(t *testing.T) {
		t.Logf("Level() = %v", tt.sl.Level())
	})

	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)
	tt = testCase[int, int]{
		name: "TestSkipList_Cap 2",
		sl:   sl,
		want: 3,
	}
	t.Run(tt.name, func(t *testing.T) {
		t.Logf("Level() = %v", tt.sl.Level())
	})
}

func TestSkipList_Cap(t *testing.T) {
	type testCase[O constraints.Ordered, T any] struct {
		name string
		sl   *SkipList[O, T]
		want int32
	}

	var sl *SkipList[int, int]
	tt := testCase[int, int]{
		name: "TestSkipList_Cap 1",
		sl:   sl,
		want: 0,
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.sl.Cap(); got != tt.want {
			t.Errorf("Cap() = %v, want %v", got, tt.want)
		}
	})

	sl = NewSkipList[int, int](10, false)
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.sl.Cap(); got != tt.want {
			t.Errorf("Cap() = %v, want %v", got, tt.want)
		}
	})

	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)
	tt = testCase[int, int]{
		name: "TestSkipList_Cap 2",
		sl:   sl,
		want: 3,
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.sl.Cap(); got != tt.want {
			t.Errorf("Cap() = %v, want %v", got, tt.want)
		}
	})

	sl.Delete(1)
	tt = testCase[int, int]{
		name: "TestSkipList_Cap 3",
		sl:   sl,
		want: 2,
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.sl.Cap(); got != tt.want {
			t.Errorf("Cap() = %v, want %v", got, tt.want)
		}
	})

	sl.Delete(5)
	tt = testCase[int, int]{
		name: "TestSkipList_Cap 4",
		sl:   sl,
		want: 2,
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := tt.sl.Cap(); got != tt.want {
			t.Errorf("Cap() = %v, want %v", got, tt.want)
		}
	})
}

func TestSkipList_Get(t *testing.T) {
	type args[O constraints.Ordered] struct {
		key O
	}
	type testCase[O constraints.Ordered, T any] struct {
		name      string
		sl        *SkipList[O, T]
		args      args[O]
		wantVal   T
		wantExist bool
	}

	var sl = NewSkipList[int, int](10, false)
	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)
	tests := []testCase[int, int]{
		{
			name:      "TestSkipList_Get 1",
			sl:        sl,
			args:      args[int]{1},
			wantVal:   1,
			wantExist: true,
		},
		{
			name:      "TestSkipList_Get 1",
			sl:        sl,
			args:      args[int]{-1},
			wantVal:   0,
			wantExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotExist := tt.sl.Get(tt.args.key)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotExist != tt.wantExist {
				t.Errorf("Get() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func TestSkipList_Put(t *testing.T) {
	type args[O constraints.Ordered, T any] struct {
		key O
		val T
	}
	type testCase[O constraints.Ordered, T any] struct {
		name string
		sl   *SkipList[O, T]
		args args[O, T]
	}

	var sl = NewSkipList[int, int](10, false)
	tt := testCase[int, int]{
		name: "TestSkipList_Put 1",
		sl:   sl,
		args: args[int, int]{1, 1},
	}
	t.Run(tt.name, func(t *testing.T) {
		tt.sl.Put(tt.args.key, tt.args.val)
	})

	tt = testCase[int, int]{
		name: "TestSkipList_Put 2",
		sl:   sl,
		args: args[int, int]{20, 20},
	}
	t.Run(tt.name, func(t *testing.T) {
		tt.sl.Put(tt.args.key, tt.args.val)
	})

	tt = testCase[int, int]{
		name: "TestSkipList_Put 3",
		sl:   sl,
		args: args[int, int]{3, 3},
	}
	t.Run(tt.name, func(t *testing.T) {
		tt.sl.Put(tt.args.key, tt.args.val)
	})

	tt = testCase[int, int]{
		name: "TestSkipList_Put 4",
		sl:   sl,
		args: args[int, int]{-3, -3},
	}
	t.Run(tt.name, func(t *testing.T) {
		tt.sl.Put(tt.args.key, tt.args.val)
	})

	tt = testCase[int, int]{
		name: "TestSkipList_Put 5",
		sl:   sl,
		args: args[int, int]{-3, 3},
	}
	t.Run(tt.name, func(t *testing.T) {
		tt.sl.Put(tt.args.key, tt.args.val)
	})
}

func TestSkipList_Delete(t *testing.T) {
	type args[O constraints.Ordered] struct {
		key O
	}
	type testCase[O constraints.Ordered, T any] struct {
		name string
		sl   *SkipList[O, T]
		args args[O]
	}

	var sl = NewSkipList[int, int](10, false)
	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)

	tests := []testCase[int, int]{
		{
			name: "TestSkipList_Delete 1",
			sl:   sl,
			args: args[int]{1},
		},
		{
			name: "TestSkipList_Delete 1",
			sl:   sl,
			args: args[int]{-1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sl.Delete(tt.args.key)
		})
	}
}

func TestSkipList_Range(t *testing.T) {
	type args[O constraints.Ordered] struct {
		start O
		end   O
	}
	type testCase[O constraints.Ordered, T any] struct {
		name string
		sl   *SkipList[O, T]
		args args[O]
		want []*KvPair[O, T]
	}

	var sl = NewSkipList[int, int](10, false)
	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)

	tests := []testCase[int, int]{
		{
			name: "TestSkipList_Range 1",
			sl:   sl,
			args: args[int]{1, 3},
			want: []*KvPair[int, int]{{1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "TestSkipList_Range 2",
			sl:   sl,
			args: args[int]{0, 5},
			want: []*KvPair[int, int]{{1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "TestSkipList_Range 3",
			sl:   sl,
			args: args[int]{3, 5},
			want: []*KvPair[int, int]{{3, 3}},
		},
		{
			name: "TestSkipList_Range 4",
			sl:   sl,
			args: args[int]{4, 5},
			want: []*KvPair[int, int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.Range(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipList_Ceil(t *testing.T) {
	type args[O constraints.Ordered] struct {
		target O
	}
	type testCase[O constraints.Ordered, T any] struct {
		name  string
		sl    *SkipList[O, T]
		args  args[O]
		want  *KvPair[O, T]
		want1 bool
	}

	var sl = NewSkipList[int, int](10, false)
	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)

	tests := []testCase[int, int]{
		{
			name:  "TestSkipList_Ceil 1",
			sl:    sl,
			args:  args[int]{0},
			want:  &KvPair[int, int]{1, 1},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 2",
			sl:    sl,
			args:  args[int]{1},
			want:  &KvPair[int, int]{1, 1},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 3",
			sl:    sl,
			args:  args[int]{2},
			want:  &KvPair[int, int]{2, 2},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 4",
			sl:    sl,
			args:  args[int]{3},
			want:  &KvPair[int, int]{3, 3},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 5",
			sl:    sl,
			args:  args[int]{4},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.sl.Ceil(tt.args.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ceil() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Ceil() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSkipList_Floor(t *testing.T) {
	type args[O constraints.Ordered] struct {
		target O
	}
	type testCase[O constraints.Ordered, T any] struct {
		name  string
		sl    *SkipList[O, T]
		args  args[O]
		want  *KvPair[O, T]
		want1 bool
	}
	var sl = NewSkipList[int, int](10, false)
	sl.Put(1, 1)
	sl.Put(2, 2)
	sl.Put(3, 3)

	tests := []testCase[int, int]{
		{
			name:  "TestSkipList_Ceil 1",
			sl:    sl,
			args:  args[int]{0},
			want:  nil,
			want1: false,
		},
		{
			name:  "TestSkipList_Ceil 2",
			sl:    sl,
			args:  args[int]{1},
			want:  &KvPair[int, int]{1, 1},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 3",
			sl:    sl,
			args:  args[int]{2},
			want:  &KvPair[int, int]{2, 2},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 4",
			sl:    sl,
			args:  args[int]{3},
			want:  &KvPair[int, int]{3, 3},
			want1: true,
		},
		{
			name:  "TestSkipList_Ceil 5",
			sl:    sl,
			args:  args[int]{4},
			want:  &KvPair[int, int]{3, 3},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.sl.Floor(tt.args.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Floor() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Floor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
