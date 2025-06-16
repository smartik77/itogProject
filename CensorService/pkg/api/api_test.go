package api

import (
	"CensorService/pkg/moderation"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"testing"
)

func TestAPI_AddForbiddenWordHandler(t *testing.T) {
	type fields struct {
		censor *moderation.CensorService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				censor: tt.fields.censor,
			}
			api.AddForbiddenWordHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestAPI_CheckHandler(t *testing.T) {
	type fields struct {
		censor *moderation.CensorService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				censor: tt.fields.censor,
			}
			api.CheckHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestAPI_Endpoints(t *testing.T) {
	type fields struct {
		censor *moderation.CensorService
	}
	type args struct {
		router *mux.Router
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				censor: tt.fields.censor,
			}
			api.Endpoints(tt.args.router)
		})
	}
}

func TestAPI_HealthHandler(t *testing.T) {
	type fields struct {
		censor *moderation.CensorService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				censor: tt.fields.censor,
			}
			api.HealthHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestNewAPI(t *testing.T) {
	type args struct {
		censor *moderation.CensorService
	}
	tests := []struct {
		name string
		args args
		want *API
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPI(tt.args.censor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
