package api

import (
	"aggregator/pkg/posts"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"testing"
)

func TestAPI_Router(t *testing.T) {
	type fields struct {
		db     *posts.DB
		router *mux.Router
	}
	tests := []struct {
		name   string
		fields fields
		want   *mux.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				db:     tt.fields.db,
				router: tt.fields.router,
			}
			if got := api.Router(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Router() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_posts(t *testing.T) {
	type fields struct {
		db     *posts.DB
		router *mux.Router
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
				db:     tt.fields.db,
				router: tt.fields.router,
			}
			api.posts(tt.args.w, tt.args.r)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		db *posts.DB
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
			if got := New(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntParam(t *testing.T) {
	type args struct {
		r            *http.Request
		name         string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIntParam(tt.args.r, tt.args.name, tt.args.defaultValue); got != tt.want {
				t.Errorf("getIntParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
