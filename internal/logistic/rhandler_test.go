package logistic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_handle(t *testing.T) {
	type args struct {
		rh raggerHandler
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HTTPError",
			args: args{rh: func(w http.ResponseWriter, r *http.Request) error {
				return Error{
					Code: http.StatusBadRequest,
					Err:  errors.New("BadRequest"),
				}
			}},
			want: "{\"code\":400,\"error\":\"BadRequest\"}\n",
		},
		{
			name: "common error",
			args: args{rh: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("common error")
			}},
			want: "{\"code\":500,\"error\":\"common error\"}\n",
		},
	}
	for _, tt := range tests {
		req := &http.Request{}
		rr := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			got := handle(tt.args.rh)
			got.ServeHTTP(rr, req)
			require.Equal(t, tt.want, rr.Body.String())

		})
	}
}
