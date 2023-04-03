package transfer

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestPrepare(t *testing.T) {
	type args struct {
		params url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    Payload
		wantErr error
	}{
		{
			name: "valid args",
			args: args{
				params: map[string][]string{
					"amount": {"2"},
					"fr":     {"x"},
					"to":     {"y"},
				},
			},
			want: Payload{
				Amount: 2,
				Fr:     "x",
				To:     "y",
			},
		},
		{
			name: "invalid args",
			args: args{
				params: map[string][]string{
					"amount": {"asdf"},
					"fr":     {"x"},
					"to":     {"y"},
				},
			},
			wantErr: fmt.Errorf("query `amount` should be an integer"),
		},
		{
			name: "invalid args",
			args: args{
				params: map[string][]string{
					"amount": {"0"},
					"fr":     {""},
					"to":     {""},
				},
			},
			wantErr: fmt.Errorf("query `amount` has default value\nquery `fr` has default value\nquery `to` has default value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Prepare(tt.args.params)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prepare() = %v, want %v", got, tt.want)
			}
		})
	}
}
