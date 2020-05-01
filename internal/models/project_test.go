package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProject_Validate(t *testing.T) {
	type fields struct {
		Title       string
		Description string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid project",
			fields: fields{
				Title:       "Test title",
				Description: "Test description",
			},
			wantErr: false,
		},
		{
			name: "invalid title",
			fields: fields{
				Title:       "",
				Description: "Test description",
			},
			wantErr: true,
		},
		{
			name: "invalid description",
			fields: fields{
				Title:       "Test title",
				Description: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
			}
			err := p.Validate()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}