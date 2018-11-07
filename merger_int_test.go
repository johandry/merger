package merger

import (
	"reflect"
	"testing"
)

type Simple struct {
	F1 int
	F2 string
}

func Test_mergeMap(t *testing.T) {
	type args struct {
		dst    interface{}
		srcMap map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Simple",
			args:    args{dst: &Simple{}, srcMap: map[string]string{"F1": "1", "F2": "one"}},
			want:    &Simple{F1: 1, F2: "one"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mergeMap(tt.args.dst, tt.args.srcMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("mergeMap() = %+v, want %+v", tt.args.dst, tt.want)
			}
		})
	}
}
