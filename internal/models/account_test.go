package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccount_Validate(t *testing.T) {
	type fields struct {
		ID        int64
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
		Token     map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid account",
			fields: fields{
				Name:     "TestName",
				Email:    "test@email.com",
				Password: "secrettestpassword",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			fields: fields{
				Name:     "",
				Email:    "test@email.com",
				Password: "secrettestpassword",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			fields: fields{
				Name:     "TestName",
				Email:    "testinvalidemail.com",
				Password: "secrettestpassword",
			},
			wantErr: true,
		},
		{
			name: "short password (len < 6)",
			fields: fields{
				Name:     "TestName",
				Email:    "testinvalidemail.com",
				Password: "pass",
			},
			wantErr: true,
		},
	}

	var err error

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				Name:     tt.fields.Name,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}

			err = a.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
