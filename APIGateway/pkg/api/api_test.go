package api

import (
	"APIGateway/pkg/models"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"testing"
)

func TestCommentsClient_AddComment(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		ctx     context.Context
		comment *models.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommentsClient{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := c.AddComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentsClient_GetComments(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		ctx    context.Context
		newsID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommentsClient{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := c.GetComments(tt.args.ctx, tt.args.newsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_AddComment(t *testing.T) {
	type fields struct {
		NewsClient     *NewsClient
		CommentsClient *CommentsClient
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
			h := &Handler{
				NewsClient:     tt.fields.NewsClient,
				CommentsClient: tt.fields.CommentsClient,
			}
			h.AddComment(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_Endpoints(t *testing.T) {
	type fields struct {
		NewsClient     *NewsClient
		CommentsClient *CommentsClient
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
			h := &Handler{
				NewsClient:     tt.fields.NewsClient,
				CommentsClient: tt.fields.CommentsClient,
			}
			h.Endpoints(tt.args.router)
		})
	}
}

func TestHandler_GetNewsDetail(t *testing.T) {
	type fields struct {
		NewsClient     *NewsClient
		CommentsClient *CommentsClient
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
			h := &Handler{
				NewsClient:     tt.fields.NewsClient,
				CommentsClient: tt.fields.CommentsClient,
			}
			h.GetNewsDetail(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetNewsList(t *testing.T) {
	type fields struct {
		NewsClient     *NewsClient
		CommentsClient *CommentsClient
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
			h := &Handler{
				NewsClient:     tt.fields.NewsClient,
				CommentsClient: tt.fields.CommentsClient,
			}
			h.GetNewsList(tt.args.w, tt.args.r)
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		newsURL     string
		commentsURL string
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.newsURL, tt.args.commentsURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsClient_GetNews(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.NewsShort
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &NewsClient{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := c.GetNews(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNews() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsClient_GetNewsDetail(t *testing.T) {
	type fields struct {
		BaseURL string
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.NewsDetail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &NewsClient{
				BaseURL: tt.fields.BaseURL,
			}
			got, err := c.GetNewsDetail(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNewsDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNewsDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
