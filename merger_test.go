package merger

import (
	"reflect"
	"testing"
)

type Simple struct {
	F1 int
	F2 string
}

type Address struct {
	City string
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

type Student struct {
	Name      string   `mapstructure:"name"`
	TextBooks []string `mapstructure:"text_books"`
}

func TestMerge(t *testing.T) {
	type args struct {
		dst    interface{}
		srcMap map[string]string
		srcs   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				dst:    &Simple{},
				srcMap: map[string]string{"F1": "10"},
				srcs: []interface{}{
					Simple{F1: 20, F2: "hello"},
					Simple{F2: "hola"},
				},
			},
			want:    &Simple{F1: 10, F2: "hello"},
			wantErr: false,
		},
		{
			name: "Person",
			args: args{
				dst:    &Person{},
				srcMap: map[string]string{"Name": "Joe", "Address.City": "LA"},
				srcs: []interface{}{
					Person{Name: "Pepe", Age: 30},
					Person{Age: 20, Address: Address{City: "San Diego"}},
				},
			},
			want:    &Person{Name: "Joe", Age: 30, Address: Address{City: "LA"}},
			wantErr: false,
		},
		{
			name: "Student",
			args: args{
				dst:    &Student{},
				srcMap: map[string]string{"text_books": "Book1, Book2, 'The Book' "},
				srcs: []interface{}{
					Student{Name: "Pepe"},
				},
			},
			want:    &Student{Name: "Pepe", TextBooks: []string{"Book1", "Book2", "The Book"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Merge(tt.args.dst, tt.args.srcMap, tt.args.srcs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("Merge() = %+v, want %+v", tt.args.dst, tt.want)
			}
		})
	}
}

func TestMergeMap(t *testing.T) {
	type args struct {
		dst     interface{}
		srcMaps []map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				dst: &Simple{},
				srcMaps: []map[string]string{
					map[string]string{"F1": "1"},
					map[string]string{"F1": "10", "F2": "ten"},
					map[string]string{"F2": "zero"},
				},
			},
			want:    &Simple{F1: 1, F2: "ten"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MergeMap(tt.args.dst, tt.args.srcMaps...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("MergeMap() = %+v, want %+v", tt.args.dst, tt.want)
			}
		})
	}
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

func TestMergeStruct(t *testing.T) {
	type args struct {
		dst  interface{}
		srcs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				dst: &Simple{},
				srcs: []interface{}{
					Simple{F1: 20},
					Simple{F2: "hola"},
					Simple{F1: 30, F2: "hello"},
				},
			},
			want:    &Simple{F1: 20, F2: "hola"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MergeStruct(tt.args.dst, tt.args.srcs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("MergeStruct() = %+v, want %+v", tt.args.dst, tt.want)
			}
		})
	}
}

func TestTransformMap(t *testing.T) {
	type args struct {
		srcMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{name: "Simple",
			args: args{srcMap: map[string]string{"F1": "20", "F2": "Text", "F3": "True"}},
			want: map[string]interface{}{"F1": "20", "F2": "Text", "F3": "True"},
		},
		{name: "Slices",
			args: args{srcMap: map[string]string{
				"Books":    "  B1,B2, B3",
				"Articles": "[A1, A2, A3]  ",
				"Items":    "  [  I1  , I2,  'Item number #3'   ,   I4  ]  ",
			}},
			want: map[string]interface{}{
				"Books":    []string{"B1", "B2", "B3"},
				"Articles": []string{"A1", "A2", "A3"},
				"Items":    []string{"I1", "I2", "Item number #3", "I4"}},
		},
		{name: "Structs",
			args: args{srcMap: map[string]string{
				"Address.City":           "New York",
				"Address.Country":        "US",
				"Parents.Address.Zip":    "32123",
				"Parents.Address.Planet": "Earth",
			}},
			want: map[string]interface{}{
				"Address": map[string]interface{}{
					"City":    "New York",
					"Country": "US",
				},
				"Parents": map[string]interface{}{
					"Address": map[string]interface{}{
						"Zip":    "32123",
						"Planet": "Earth",
					},
				},
			},
		},
		{name: "Mixed",
			args: args{srcMap: map[string]string{"IP": "192.168.1.0", "DNS.Servers": "[192.168.0.1, 192.168.0.2, 192.168.0.3]"}},
			want: map[string]interface{}{
				"IP": "192.168.1.0",
				"DNS": map[string]interface{}{
					"Servers": []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformMap(tt.args.srcMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
