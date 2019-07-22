package merger_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/johandry/merger"
)

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
				"OneItem":  "[Item]",
				"Books":    "  B1,B2, B3",
				"Articles": "[A1, A2, A3]  ",
				"Items":    "  [  I1  , I2,  'Item number #3'   ,   I4  ]  ",
			}},
			want: map[string]interface{}{
				"OneItem":  []string{"Item"},
				"Books":    []string{"B1", "B2", "B3"},
				"Articles": []string{"A1", "A2", "A3"},
				"Items":    []string{"I1", "I2", "Item number #3", "I4"},
			},
		},
		{name: "Structs",
			args: args{srcMap: map[string]string{
				"Address__City":            "New York",
				"Address__Country":         "US",
				"Parents__Address__Zip":    "32123",
				"Parents__Address__Planet": "Earth",
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
			args: args{srcMap: map[string]string{"IP": "192.168.1.0", "DNS__Servers": "[192.168.0.1, 192.168.0.2, 192.168.0.3]"}},
			want: map[string]interface{}{
				"IP": "192.168.1.0",
				"DNS": map[string]interface{}{
					"Servers": []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"},
				},
			},
		},
		{name: "Structs-JSON",
			args: args{srcMap: map[string]string{
				"Address": `{"city": "New York", "country": "US"}`,
				"Parents": `{"address": {"zip": "32123", "planet": "Earth"}}`,
			}},
			want: map[string]interface{}{
				"Address": map[string]interface{}{
					"city":    "New York",
					"country": "US",
				},
				"Parents": map[string]interface{}{
					"address": map[string]interface{}{
						"zip":    "32123",
						"planet": "Earth",
					},
				},
			},
		},
		{name: "Mixed-JSON",
			args: args{srcMap: map[string]string{
				"IP":                    "192.168.1.0",
				"DNS__Servers":          "[192.168.0.1, 192.168.0.2, 192.168.0.3]",
				"Parents__Address__Zip": "32123",
				"Parents":               `{"Address": {"Planet": "Earth"}}`,
			}},
			want: map[string]interface{}{
				"IP": "192.168.1.0",
				"DNS": map[string]interface{}{
					"Servers": []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"},
				},
				"Parents": map[string]interface{}{
					"Address": map[string]interface{}{
						"Zip":    "32123",
						"Planet": "Earth",
					},
				},
			},
		},
	}
	// assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := merger.TransformMap(tt.args.srcMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformMap() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

type SimpleWithIgnore struct {
	F1 int    `json:"age" mapstructure:"-"`
	F2 string `yaml:"name" mapstructure:"-"`
}

type Movies struct {
	DBName string   `json:"movie_database_name"`
	Movies []*Movie `json:"movies"`
}

type Movie struct {
	Title     string                 `json:"title"`
	Year      int                    `json:"year"`
	Comment   string                 `json:"-"`
	Actors    map[string]MoviePerson `json:"actors"`
	Genres    []string               `json:"genres"`
	Release   map[string]int         `json:"release_year_per_country"`
	Directors []MoviePerson          `json:"directors"`
}

type MoviePerson struct {
	Name string `json:"full_name"`
	Age  int    `json:"age"`
}

var movie01 = Movie{
	Title:   "Seven Samurai",
	Year:    1954,
	Comment: "Good movie, but this comment goes to the void",
	Actors: map[string]MoviePerson{
		"Kikuchiyo": MoviePerson{
			Name: "Toshirô Mifune",
		},
		"Kambei Shimada": MoviePerson{
			Name: "Takashi Shimura",
		},
	},
	Genres: []string{"Adventure", "Drama"},
	Release: map[string]int{
		"JP": 1954,
		"US": 1956,
	},
	Directors: []MoviePerson{
		MoviePerson{Name: "Akira Kurosawa"},
	},
}
var movie01Map = map[string]string{
	"title":                             "Seven Samurai",
	"year":                              "1954",
	"actors__kikuchiyo__full_name":      "Toshirô Mifune",
	"actors__kikuchiyo__age":            "0",
	"actors__kambei_shimada__full_name": "Takashi Shimura",
	"actors__kambei_shimada__age":       "0",
	"genres":                            "[Adventure, Drama]",
	"release_year_per_country__JP":      "1954",
	"release_year_per_country__US":      "1956",
}

var movie02 = Movie{
	Title:   "The Godfather",
	Year:    1972,
	Comment: "Good movie, but this comment will be ignored",
	Actors: map[string]MoviePerson{
		"Don Vito Corleone": MoviePerson{
			Name: "Marlon Brando",
		},
		"Michael Corleone": MoviePerson{
			Name: "Al Pacino",
			Age:  time.Now().Year() - 1940,
		},
	},
	Genres: []string{"Crime", "Drama"},
	Release: map[string]int{
		"HK": 1973,
		"US": 1972,
	},
	Directors: []MoviePerson{
		MoviePerson{Name: "Francis Ford Coppola"},
	},
}
var movie02Map = map[string]string{
	"title":                                "The Godfather",
	"year":                                 "1972",
	"actors__don_vito_corleone__full_name": "Marlon Brando",
	"actors__don_vito_corleone__age":       "0",
	"actors__michael_corleone__full_name":  "Al Pacino",
	"actors__michael_corleone__age":        "79",
	"genres":                               "[Crime, Drama]",
	"release_year_per_country__HK":         "1973",
	"release_year_per_country__US":         "1972",
}

func TestTransformToMap(t *testing.T) {
	i := 10

	tests := []struct {
		name    string
		v       interface{}
		tags    []string
		want    map[string]string
		wantErr bool
	}{
		{name: "nil",
			v:       nil,
			want:    map[string]string{},
			wantErr: false,
		},
		{name: "no struct (int)",
			v:       10,
			want:    map[string]string{},
			wantErr: true,
		},
		{name: "no struct pointer",
			v:       Simple{1, "foo"},
			want:    map[string]string{},
			wantErr: true,
		},
		{name: "simple with ignored fields",
			v: &SimpleWithIgnore{
				F1: 30,
				F2: "John",
			},
			tags:    []string{"mapstructure", "json", "yaml"},
			want:    map[string]string{},
			wantErr: false,
		},
		{name: "simple",
			v: &Simple{
				F1: 123,
				F2: "something",
			},
			want: map[string]string{
				"f1": "123",
				"f2": "something",
			},
			wantErr: false,
		},
		{name: "complex",
			v: &Person{
				Name:    "Pepe",
				Age:     20,
				Address: Address{City: "San Diego", Country: "US"},
				Phones: map[string]Phone{
					"home":   Phone{Number: "858-123-4567", Available: true},
					"mobile": Phone{Number: "858-987-6543"},
				},
			},
			want: map[string]string{
				"name":                      "Pepe",
				"age":                       "20",
				"address__city":             "San Diego",
				"address__country":          "US",
				"phones__home__number":      "858-123-4567",
				"phones__home__available":   "true",
				"phones__mobile__number":    "858-987-6543",
				"phones__mobile__available": "false",
			},
			wantErr: false,
		},
		{name: "movies",
			v: &Movies{
				DBName: "DBTest",
				Movies: []*Movie{
					&movie01,
					&movie02,
				},
			},
			want: map[string]string{
				"movie_database_name": "DBTest",
			},
			wantErr: false,
		},
		{name: "movie01",
			v:       &movie01,
			want:    movie01Map,
			wantErr: false,
		},
		{name: "movie02",
			v:       &movie02,
			want:    movie02Map,
			wantErr: false,
		},
		{name: "empty list",
			v: &struct {
				Empty    []int
				NotEmpty []string
			}{
				Empty:    []int{},
				NotEmpty: []string{"a", "b"},
			},
			want: map[string]string{
				"empty":    "[]",
				"notempty": "[a, b]",
			},
			wantErr: false,
		},
		{name: "unexported elements & pointer",
			v: &struct {
				private    int
				PublicList []string
				Ptr        *int
			}{
				private:    10,
				PublicList: []string{"a", "b"},
				Ptr:        &i,
			},
			want: map[string]string{
				"publiclist": "[a, b]",
				"ptr":        "10",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := merger.TransformToMap(tt.v, tt.tags...)
			t.Logf("Test %q returns %v", tt.name, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
