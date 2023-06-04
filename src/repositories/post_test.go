package repositories

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igorsfreitas/devbook-api/src/models"
)

func Test_NewPostRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}

	tests := []struct {
		name string
		args args
		want *Posts
	}{
		{
			name: "ssccess creating a new post repository",
			args: args{
				db: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostRepository(tt.args.db); got == nil {
				t.Errorf("[Test_NewPostRepository] NewPostRepository() = %v", got)
			}
		})
	}
}

func Test_CreatePost(t *testing.T) {
	type args struct {
		post models.Post
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       uint64
		wantErr    bool
	}{
		{
			name: "success creating a new post",
			args: args{
				post: models.Post{
					Title:    "title",
					Content:  "content",
					AuthorID: 1,
				},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id").
					ExpectQuery().
					WithArgs("title", "content", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
			},
			want: 1,
		},
		{
			name: "error creating a new post - prepare error",
			args: args{
				post: models.Post{
					Title:    "title",
					Content:  "content",
					AuthorID: 1,
				},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error creating a new post - query error",
			args: args{
				post: models.Post{
					Title:    "title",
					Content:  "content",
					AuthorID: 1,
				},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id").
					ExpectQuery().
					WithArgs("title", "content", 1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := p.Create(tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_CreatePost] Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_CreatePost] Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPost(t *testing.T) {
	type args struct {
		postID uint64
	}

	success := models.Post{
		ID:         1,
		Title:      "title",
		Content:    "content",
		AuthorID:   1,
		Likes:      1,
		AuthorNick: "nick",
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.Post
		wantErr    bool
	}{
		{
			name: "success getting a post",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick
							FROM posts p
							INNER JOIN users u
								ON u.id = p.author_id
							WHERE p.id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "content", "likes", "created_at", "nick"}).
						AddRow(success.ID, success.AuthorID, success.Title, success.Content, success.Likes, success.CreatedAt, success.AuthorNick))
			},
			want: success,
		},
		{
			name: "error getting a post - query error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick
							FROM posts p
							INNER JOIN users u
								ON u.id = p.author_id
							WHERE p.id = $1
					`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error getting a post - scan error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick
							FROM posts p
							INNER JOIN users u
								ON u.id = p.author_id
							WHERE p.id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
						AddRow(nil, nil, nil, nil, nil, time.Now(), nil))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := p.GetPost(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_GetPost] Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_GetPost] Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetPosts(t *testing.T) {
	type args struct {
		userID uint64
	}

	success := models.Post{
		ID:         1,
		Title:      "title",
		Content:    "content",
		AuthorID:   1,
		Likes:      1,
		AuthorNick: "nick",
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       []models.Post
		wantErr    bool
	}{
		{
			name: "success getting posts",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT DISTINCT p.*, u.nick FROM posts p
						INNER JOIN users u ON u.id = p.author_id
						INNER JOIN followers f ON p.author_id = f.user_id
						WHERE u.id = $1 OR f.follower_id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "content", "likes", "created_at", "nick"}).
						AddRow(success.ID, success.AuthorID, success.Title, success.Content, success.Likes, success.CreatedAt, success.AuthorNick))
			},
			want: []models.Post{success},
		},
		{
			name: "error getting posts - query error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT DISTINCT p.*, u.nick FROM posts p
						INNER JOIN users u ON u.id = p.author_id
						INNER JOIN followers f ON p.author_id = f.user_id
						WHERE u.id = $1 OR f.follower_id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error getting posts - scan error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT DISTINCT p.*, u.nick FROM posts p
						INNER JOIN users u ON u.id = p.author_id
						INNER JOIN followers f ON p.author_id = f.user_id
						WHERE u.id = $1 OR f.follower_id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
						AddRow(nil, nil, nil, nil, nil, time.Now(), nil))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := p.GetPosts(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_GetPosts] Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_GetPosts] Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdatePost(t *testing.T) {
	type args struct {
		postID uint64
		post   models.Post
	}

	success := models.Post{
		ID:         1,
		Title:      "title",
		Content:    "content",
		AuthorID:   1,
		Likes:      1,
		AuthorNick: "nick",
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.Post
		wantErr    bool
	}{
		{
			name: "success updating post",
			args: args{
				postID: 1,
				post:   success,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET title = $1, content = $2 WHERE id = $3
					`).
					ExpectExec().
					WithArgs(success.Title, success.Content, success.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: success,
		},
		{
			name: "error updating post - prepare error",
			args: args{
				postID: 1,
				post:   success,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET title = $1, content = $2 WHERE id = $3
					`).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error updating post - exec error",
			args: args{
				postID: 1,
				post:   success,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET title = $1, content = $2 WHERE id = $3
					`).
					ExpectExec().
					WithArgs(success.Title, success.Content, success.ID).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = p.Update(tt.args.postID, tt.args.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_UpdatePost] Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_DeletePost(t *testing.T) {
	type args struct {
		postID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success deleting post",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						DELETE FROM posts WHERE id = $1
					`).
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error deleting post - prepare error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						DELETE FROM posts WHERE id = $1
					`).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error deleting post - exec error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						DELETE FROM posts WHERE id = $1
					`).
					ExpectExec().
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = p.Delete(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_DeletePost] Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_GetPostsByUser(t *testing.T) {
	type args struct {
		userID uint64
	}

	success := models.Post{
		ID:         1,
		Title:      "title",
		Content:    "content",
		AuthorID:   1,
		Likes:      1,
		AuthorNick: "nick",
		CreatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       []models.Post
		wantErr    bool
	}{
		{
			name: "success getting posts by user",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick FROM posts p
						JOIN users u ON u.id = p.author_id
						WHERE u.id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "content", "likes", "created_at", "nick"}).
						AddRow(success.ID, success.AuthorID, success.Title, success.Content, success.Likes, success.CreatedAt, success.AuthorNick))
			},
			want: []models.Post{success},
		},
		{
			name: "error getting posts by user - query error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick FROM posts p
						JOIN users u ON u.id = p.author_id
						WHERE u.id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error getting posts by user - scan error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(`
						SELECT p.*, u.nick FROM posts p
						JOIN users u ON u.id = p.author_id
						WHERE u.id = $1
						ORDER BY 1 DESC
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "author_id", "likes", "created_at", "nick"}).
						AddRow(nil, nil, nil, nil, nil, time.Now(), nil))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := p.GetPostsByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_GetPostsByUser] Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_GetPostsByUser] Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Like(t *testing.T) {
	type args struct {
		postID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success liking post",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET likes = likes + 1 WHERE id = $1
					`).
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error liking post - prepare error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET likes = likes + 1 WHERE id = $1
					`).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error liking post - exec error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
						UPDATE posts SET likes = likes + 1 WHERE id = $1
					`).
					ExpectExec().
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = p.Like(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_Like] Like() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_Unlike(t *testing.T) {
	type args struct {
		postID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success unliking post",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
					UPDATE posts 
						SET likes = 
							CASE WHEN likes > 0 
								THEN likes - 1 
								ELSE 0 
							END 
						WHERE id = $1`).
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error unliking post - prepare error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
					UPDATE posts 
						SET likes = 
							CASE WHEN likes > 0 
								THEN likes - 1 
								ELSE 0 
							END 
						WHERE id = $1`).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "error unliking post - exec error",
			args: args{
				postID: 1,
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectPrepare(`
					UPDATE posts 
						SET likes = 
							CASE WHEN likes > 0 
								THEN likes - 1 
								ELSE 0 
							END 
						WHERE id = $1`).
					ExpectExec().
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			p := Posts{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = p.Unlike(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("[Test_Unlike] Unlike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
