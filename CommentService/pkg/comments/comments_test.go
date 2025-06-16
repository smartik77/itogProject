package comments

import (
	"CommentService/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
	"reflect"
	"testing"
)

func TestCommentRepository_AddComment(t *testing.T) {
	type fields struct {
		Conn *pgx.Conn
	}
	type args struct {
		ctx     context.Context
		comment models.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommentRepository{
				Conn: tt.fields.Conn,
			}
			got, err := r.AddComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentRepository_GetCommentContent(t *testing.T) {
	type fields struct {
		Conn *pgx.Conn
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommentRepository{
				Conn: tt.fields.Conn,
			}
			got, err := r.GetCommentContent(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCommentContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentRepository_GetCommentsByNewsID(t *testing.T) {
	type fields struct {
		Conn *pgx.Conn
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
			r := &CommentRepository{
				Conn: tt.fields.Conn,
			}
			got, err := r.GetCommentsByNewsID(tt.args.ctx, tt.args.newsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByNewsID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommentsByNewsID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentRepository_UpdateModerationStatus(t *testing.T) {
	type fields struct {
		Conn *pgx.Conn
	}
	type args struct {
		ctx    context.Context
		id     int
		status string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommentRepository{
				Conn: tt.fields.Conn,
			}
			if err := r.UpdateModerationStatus(tt.args.ctx, tt.args.id, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateModerationStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
