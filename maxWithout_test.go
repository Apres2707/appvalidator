package appvalidator

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func Test_isLessOrEqual(t *testing.T) {
	now := time.Now()
	dateTime := now.Add(-1 * time.Duration(now.Nanosecond()))

	type args struct {
		value          reflect.Value
		maxValueString string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "less uint64",
			args: args{
				value:          reflect.ValueOf(uint64(9)),
				maxValueString: "10",
			},
			want: true,
		},
		{
			name: "equal uint64",
			args: args{
				value:          reflect.ValueOf(10),
				maxValueString: "10",
			},
			want: true,
		},
		{
			name: "great uint64",
			args: args{
				value:          reflect.ValueOf(11),
				maxValueString: "10",
			},
			want: false,
		},
		{
			name: "less string",
			args: args{
				value:          reflect.ValueOf("a"),
				maxValueString: "2",
			},
			want: true,
		},
		{
			name: "equal string",
			args: args{
				value:          reflect.ValueOf("ab"),
				maxValueString: "2",
			},
			want: true,
		},
		{
			name: "great string",
			args: args{
				value:          reflect.ValueOf("abc"),
				maxValueString: "2",
			},
			want: false,
		},
		{
			name: "less slice",
			args: args{
				value:          reflect.ValueOf([]int64{1}),
				maxValueString: "2",
			},
			want: true,
		},
		{
			name: "equal slice",
			args: args{
				value:          reflect.ValueOf([]int64{1, 2}),
				maxValueString: "2",
			},
			want: true,
		},
		{
			name: "great slice",
			args: args{
				value:          reflect.ValueOf([]int64{1, 2, 3}),
				maxValueString: "2",
			},
			want: false,
		},
		{
			name: "empty slice",
			args: args{
				value:          reflect.ValueOf([]int64{}),
				maxValueString: "2",
			},
			want: true,
		},
		{
			name: "less date",
			args: args{
				value:          reflect.ValueOf(dateTime),
				maxValueString: dateTime.Add(time.Hour).Format("2006-01-02T15:04:05-07:00"),
			},
			want: true,
		},
		{
			name: "equal date",
			args: args{
				value:          reflect.ValueOf(dateTime),
				maxValueString: dateTime.Format("2006-01-02T15:04:05-07:00"),
			},
			want: true,
		},
		{
			name: "great date",
			args: args{
				value:          reflect.ValueOf(dateTime),
				maxValueString: dateTime.Add(-1 * time.Hour).Format("2006-01-02T15:04:05-07:00"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLessOrEqual(tt.args.value, tt.args.maxValueString); got != tt.want {
				t.Errorf("isLessOrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPropertyNotEmpty(t *testing.T) {
	testString := "abc"

	type args struct {
		val      reflect.Value
		pathList []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string not empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop string
					}{
						Prop: "abc",
					},
				),
				pathList: []string{"Prop"},
			},
			want: true,
		},
		{
			name: "string empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop string
					}{
						Prop: "",
					},
				),
				pathList: []string{"Prop"},
			},
			want: false,
		},
		{
			name: "ptr string nil",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop *string
					}{
						Prop: nil,
					},
				),
				pathList: []string{"Prop"},
			},
			want: false,
		},
		{
			name: "ptr string not empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop *string
					}{
						Prop: &testString,
					},
				),
				pathList: []string{"Prop"},
			},
			want: true,
		},
		{
			name: "slice not empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop []string
					}{
						Prop: []string{"abc"},
					},
				),
				pathList: []string{"Prop"},
			},
			want: true,
		},
		{
			name: "slice empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop []string
					}{
						Prop: []string{},
					},
				),
				pathList: []string{"Prop"},
			},
			want: false,
		},
		{
			name: "slice nil",
			args: args{
				val: reflect.ValueOf(
					struct {
						Prop []string
					}{
						Prop: nil,
					},
				),
				pathList: []string{"Prop"},
			},
			want: false,
		},
		{
			name: "nested slice not empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						NestedProp struct {
							Prop []string
						}
					}{
						NestedProp: struct {
							Prop []string
						}{
							Prop: []string{"abc"},
						},
					},
				),
				pathList: []string{"NestedProp", "Prop"},
			},
			want: true,
		},
		{
			name: "nested slice empty",
			args: args{
				val: reflect.ValueOf(
					struct {
						NestedProp struct {
							Prop []string
						}
					}{
						NestedProp: struct {
							Prop []string
						}{
							Prop: []string{},
						},
					},
				),
				pathList: []string{"NestedProp", "Prop"},
			},
			want: false,
		},
		{
			name: "nested slice nil",
			args: args{
				val: reflect.ValueOf(
					struct {
						NestedProp struct {
							Prop []string
						}
					}{
						NestedProp: struct {
							Prop []string
						}{
							Prop: nil,
						},
					},
				),
				pathList: []string{"NestedProp", "Prop"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPropertyNotEmpty(tt.args.val, tt.args.pathList); got != tt.want {
				t.Errorf("isPropertyNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Validate(t *testing.T) {
	limit := uint64(20)

	type args struct {
		data testData
		tag  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "great than max with not empty dependency",
			args: args{
				data: testData{
					Limit: &limit,
					Filters: filters{
						StringF: "a",
					},
				},
				tag: `validate:"max_without=Filters.StringF 10"`,
			},
			wantErr: false,
		},
		{
			name: "great than max with empty dependency",
			args: args{
				data: testData{
					Limit: &limit,
					Filters: filters{
						StringF: "",
					},
				},
				tag: `validate:"max_without=Filters.StringF 10"`,
			},
			wantErr: true,
		},
		{
			name: "great than max with not empty nested dependency",
			args: args{
				data: testData{
					Limit: &limit,
					Filters: filters{
						NestedFilters: nestedFilters{
							SliceF: []string{"abc"},
						},
					},
				},
				tag: `validate:"max_without=Filters.NestedFilters.SliceF 10"`,
			},
			wantErr: false,
		},
		{
			name: "great than max with empty nested dependency",
			args: args{
				data: testData{
					Limit: &limit,
					Filters: filters{
						NestedFilters: nestedFilters{
							SliceF: []string{},
						},
					},
				},
				tag: `validate:"max_without=Filters.NestedFilters.SliceF 10"`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		v := validator.New()
		err := WithCustom(v)
		if err != nil {
			t.Errorf("WithCustom(): %s", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			taggedData := buildWithLimitValidationTag(tt.args.data, tt.args.tag)
			if err != nil {
				t.Errorf("tag test data: %s", err)
			}

			err := v.Struct(taggedData)
			if (err != nil) != tt.wantErr {
				t.Errorf("An error is expected %v but got %v.", tt.wantErr, err)
			}
		})
	}
}

type testData struct {
	Limit   *uint64
	Filters filters
}

type filters struct {
	StringF       string
	NestedFilters nestedFilters
}

type nestedFilters struct {
	SliceF []string
}

func buildWithLimitValidationTag(rawData testData, tag string) interface{} {
	value := reflect.ValueOf(rawData)
	t := value.Type()
	sf := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Tag)
		sf = append(sf, t.Field(i))
		if t.Field(i).Name == "Limit" {
			sf[i].Tag = reflect.StructTag(tag)
		}
	}

	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)

	return newValue.Interface()
}
