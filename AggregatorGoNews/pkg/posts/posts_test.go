package posts

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestDB_InsertPost(t *testing.T) {
	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		post *Post
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
			db := &DB{
				pool: tt.fields.pool,
			}
			if err := db.InsertPost(tt.args.post); (err != nil) != tt.wantErr {
				t.Errorf("InsertPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_SearchPosts(t *testing.T) {
	type fields struct {
		pool *pgxpool.Pool
	}
	type args struct {
		search string
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Post
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				pool: tt.fields.pool,
			}
			got, got1, err := db.SearchPosts(tt.args.search, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchPosts() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SearchPosts() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
