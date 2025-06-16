package api

import (
	"CommentService/pkg/comments"
	"github.com/jackc/pgx/v5"
	"net/http"
	"reflect"
	"testing"
)

func TestCommentHandler_AddCommentHandler(t *testing.T) {
	type fields struct {
		Repo  *comments.CommentRepository
		Queue chan int
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
			h := &CommentHandler{
				Repo:  tt.fields.Repo,
				Queue: tt.fields.Queue,
			}
			h.AddCommentHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestCommentHandler_GetCommentsHandler(t *testing.T) {
	type fields struct {
		Repo  *comments.CommentRepository
		Queue chan int
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
			h := &CommentHandler{
				Repo:  tt.fields.Repo,
				Queue: tt.fields.Queue,
			}
			h.GetCommentsHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestCommentHandler_ModerateCommentHandler(t *testing.T) {
	type fields struct {
		Repo  *comments.CommentRepository
		Queue chan int
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
			h := &CommentHandler{
				Repo:  tt.fields.Repo,
				Queue: tt.fields.Queue,
			}
			h.ModerateCommentHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestEndpoints(t *testing.T) {
	type args struct {
		conn *pgx.Conn
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Endpoints(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCommentHandler(t *testing.T) {
	type args struct {
		conn *pgx.Conn
	}
	tests := []struct {
		name string
		args args
		want *CommentHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommentHandler(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startModerationWorker(t *testing.T) {
	type args struct {
		repo  *comments.CommentRepository
		queue chan int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startModerationWorker(tt.args.repo, tt.args.queue)
		})
	}
}
