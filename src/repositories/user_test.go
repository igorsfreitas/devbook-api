package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/igorsfreitas/devbook-api/src/models"
)

func Test_NewUserRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}

	tests := []struct {
		name string
		args args
		want *Users
	}{
		{
			name: "success creating a new user repository",
			args: args{
				db: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.db); got == nil {
				t.Errorf("[Test_NewUserRepository] NewUserRepository() = %v", got)
			}
		})
	}
}

func Test_Create(t *testing.T) {

	type args struct {
		user models.User
	}

	userSuccess := models.User{
		Name:     "Test",
		Nick:     "testnick",
		Email:    "nick@test.com.br",
		Password: "123456",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       uint64
		wantErr    bool
	}{
		{
			name: "success creating a user",
			args: args{
				user: userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into users (name, nick, email, password) values ($1, $2, $3, $4) returning id").
					ExpectQuery().
					WithArgs(userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want: 1,
		},
		{
			name: "error creating a user - query error",
			args: args{
				user: userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into users (name, nick, email, password) values ($1, $2, $3, $4) returning id").
					ExpectQuery().
					WithArgs(userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.Password).
					WillReturnError(errors.New("error creating a user"))
			},
			wantErr: true,
		},
		{
			name: "error creating a user - prepare error",
			args: args{
				user: userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into users (name, nick, email, password) values ($1, $2, $3, $4) returning id").
					WillReturnError(errors.New("error creating a user"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_CreateUser] Users.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Find(t *testing.T) {
	type args struct {
		nickOrName string
	}

	nickOrName := fmt.Sprintf("%%%s%%", "testnick")

	userSuccess := models.User{
		ID:        1,
		Name:      "Test",
		Nick:      "testnick",
		Email:     "nick@test.com.br",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       []models.User
		wantErr    bool
	}{
		{
			name: "success finding a users",
			args: args{
				nickOrName: userSuccess.Nick,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, name, nick, email, created_at from users where lower(name) like $1 or lower(nick) like $2").
					WithArgs(nickOrName, nickOrName).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(1, userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.CreatedAt))
			},
			want: []models.User{userSuccess},
		},
		{
			name: "error finding a users - query error",
			args: args{
				nickOrName: userSuccess.Nick,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("select id, name, nick, email, created_at from users where nick = $1").
					ExpectQuery().
					WithArgs(userSuccess.Nick).
					WillReturnError(errors.New("error finding a users"))
			},
			wantErr: true,
		},
		{
			name: "error finding a users - scan error",
			args: args{
				nickOrName: userSuccess.Nick,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, name, nick, email, created_at from users where lower(name) like $1 or lower(nick) like $2").
					WithArgs(nickOrName, nickOrName).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(nil, nil, nil, nil, nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.Find(tt.args.nickOrName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_Find] Users.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindByID(t *testing.T) {
	type args struct {
		id uint64
	}

	userSuccess := models.User{
		ID:        1,
		Name:      "Test",
		Nick:      "testnick",
		Email:     "",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.User
		wantErr    bool
	}{
		{
			name: "success finding a user by id",
			args: args{
				id: userSuccess.ID,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, name, nick, email, created_at from users where id = $1").
					WithArgs(userSuccess.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(userSuccess.ID, userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.CreatedAt))
			},
			want: userSuccess,
		},
		{
			name: "error finding a user by id - query error",
			args: args{
				id: userSuccess.ID,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("select id, name, nick, email, created_at from users where id = $1").
					ExpectQuery().
					WithArgs(userSuccess.ID).
					WillReturnError(errors.New("error finding a user by id"))
			},
			wantErr: true,
		},
		{
			name: "error finding a user by id - scan error",
			args: args{
				id: userSuccess.ID,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, name, nick, email, created_at from users where id = $1").
					WithArgs(userSuccess.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(nil, nil, nil, nil, nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_FindByID] Users.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindByEmail(t *testing.T) {
	type args struct {
		email string
	}

	userSuccess := models.User{
		ID:       1,
		Email:    "",
		Password: "123456",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.User
		wantErr    bool
	}{
		{
			name: "success finding a user by email",
			args: args{
				email: userSuccess.Email,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, password from users where email = $1").
					WithArgs(userSuccess.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow(userSuccess.ID, userSuccess.Password))
			},
			want: userSuccess,
		},
		{
			name: "error finding a user by email - query error",
			args: args{
				email: userSuccess.Email,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("select id, password from users where email = $1").
					ExpectQuery().
					WithArgs(userSuccess.Email).
					WillReturnError(errors.New("error finding a user by email"))
			},
			wantErr: true,
		},
		{
			name: "error finding a user by email - scan error",
			args: args{
				email: userSuccess.Email,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery("select id, password from users where email = $1").
					WithArgs(userSuccess.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).
						AddRow(nil, nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.FindByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[Test_FindByEmail] Users.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	type args struct {
		userID uint64
		user   models.User
	}

	userSuccess := models.User{
		ID:    1,
		Name:  "Test",
		Nick:  "test",
		Email: "",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       models.User
		wantErr    bool
	}{
		{
			name: "success updating a user",
			args: args{
				userID: userSuccess.ID,
				user:   userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("update users set name = $1, nick = $2, email = $3 where id = $4").
					ExpectExec().
					WithArgs(userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: userSuccess,
		},
		{
			name: "error updating a user - exec error",
			args: args{
				userID: userSuccess.ID,
				user:   userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("update users set name = $1, nick = $2, email = $3 where id = $4").
					ExpectExec().
					WithArgs(userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.ID).
					WillReturnError(errors.New("error updating a user"))
			},
			wantErr: true,
		},
		{
			name: "error updating a user - prepare error",
			args: args{
				userID: userSuccess.ID,
				user:   userSuccess,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("update users set name = $1, nick = $2, email = $3 where id = $4").
					WillReturnError(errors.New("error updating a user"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = u.Update(tt.args.userID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	type args struct {
		userID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success deleting a user",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from users where id = $1").
					ExpectExec().
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error deleting a user - exec error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from users where id = $1").
					ExpectExec().
					WithArgs(1).
					WillReturnError(errors.New("error deleting a user"))
			},
			wantErr: true,
		},
		{
			name: "error deleting a user - prepare error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from users where id = $1").
					WillReturnError(errors.New("error deleting a user"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = u.Delete(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_FollowUser(t *testing.T) {
	type args struct {
		userID   uint64
		followID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success following a user",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into followers (user_id, follower_id) values ($1, $2) ON CONFLICT DO NOTHING").
					ExpectExec().
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error following a user - exec error",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into followers (user_id, follower_id) values ($1, $2) ON CONFLICT DO NOTHING").
					ExpectExec().
					WithArgs(1, 2).
					WillReturnError(errors.New("error following a user"))
			},
			wantErr: true,
		},
		{
			name: "error following a user - prepare error",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("insert into followers (user_id, follower_id) values ($1, $2) ON CONFLICT DO NOTHING").
					WillReturnError(errors.New("error following a user"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = u.FollowUser(tt.args.userID, tt.args.followID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.FollowUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_UnfollowUser(t *testing.T) {
	type args struct {
		userID   uint64
		followID uint64
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success unfollowing a user",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from followers where user_id = $1 and follower_id = $2").
					ExpectExec().
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error unfollowing a user - exec error",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from followers where user_id = $1 and follower_id = $2").
					ExpectExec().
					WithArgs(1, 2).
					WillReturnError(errors.New("error unfollowing a user"))
			},
			wantErr: true,
		},
		{
			name: "error unfollowing a user - prepare error",
			args: args{
				userID:   1,
				followID: 2,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare("delete from followers where user_id = $1 and follower_id = $2").
					WillReturnError(errors.New("error unfollowing a user"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err = u.UnfollowUser(tt.args.userID, tt.args.followID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.UnfollowUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_FindFollowers(t *testing.T) {
	type args struct {
		userID uint64
	}

	userSuccess := models.User{
		ID:        1,
		Name:      "Name",
		Nick:      "Nick",
		Email:     "Email",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       []models.User
		wantErr    bool
	}{
		{
			name: "success finding followers",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
					select u.id, u.name, u.nick, u.email, u.created_at 
						from users u
						inner join followers f on u.id = f.follower_id
						where f.user_id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(userSuccess.ID, userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.CreatedAt))
			},
			want: []models.User{userSuccess},
		},
		{
			name: "error finding followers - query error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
					select u.id, u.name, u.nick, u.email, u.created_at 
						from users u
						inner join followers f on u.id = f.follower_id
						where f.user_id = $1
				`).
					WithArgs(1).
					WillReturnError(errors.New("error finding followers"))
			},
			wantErr: true,
		},
		{
			name: "error finding followers - scan error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
						select u.id, u.name, u.nick, u.email, u.created_at 
							from users u
							inner join followers f on u.id = f.follower_id
							where f.user_id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(nil, nil, nil, nil, nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.FindFollowers(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.FindFollowers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.FindFollowers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindFollowing(t *testing.T) {
	type args struct {
		userID uint64
	}

	userSuccess := models.User{
		ID:        1,
		Name:      "Name",
		Nick:      "Nick",
		Email:     "Email",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       []models.User
		wantErr    bool
	}{
		{
			name: "success finding following",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
					select u.id, u.name, u.nick, u.email, u.created_at 
						from users u
						inner join followers f on u.id = f.user_id
						where f.follower_id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(userSuccess.ID, userSuccess.Name, userSuccess.Nick, userSuccess.Email, userSuccess.CreatedAt))
			},
			want: []models.User{userSuccess},
		},
		{
			name: "error finding following - query error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
					select u.id, u.name, u.nick, u.email, u.created_at 
						from users u
						inner join followers f on u.id = f.user_id
						where f.follower_id = $1
				`).
					WithArgs(1).
					WillReturnError(errors.New("error finding following"))
			},
			wantErr: true,
		},
		{
			name: "error finding following - scan error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
						select u.id, u.name, u.nick, u.email, u.created_at 
							from users u
							inner join followers f on u.id = f.user_id
							where f.follower_id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "nick", "email", "created_at"}).
						AddRow(nil, nil, nil, nil, nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.FindFollowing(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.FindFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.FindFollowing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UpdatePassword(t *testing.T) {
	type args struct {
		userID   uint64
		password string
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success updating password",
			args: args{
				userID:   1,
				password: "password",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare(`
						update users set password = $1 where id = $2
					`).
					ExpectExec().
					WithArgs("password", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error updating password - exec error",
			args: args{
				userID:   1,
				password: "password",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare(`
						update users set password = $1 where id = $2
					`).
					ExpectExec().
					WithArgs("password", 1).
					WillReturnError(errors.New("error updating password"))
			},
			wantErr: true,
		},
		{
			name: "error updating password - prepare error",
			args: args{
				userID:   1,
				password: "password",
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectPrepare(`
						update users set password = $1 where id = $2
					`).
					WillReturnError(errors.New("error updating password"))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			if err := u.UpdatePassword(tt.args.userID, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Users.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_GetUserPassword(t *testing.T) {
	type args struct {
		userID uint64
	}

	userSuccess := models.User{
		ID:       1,
		Password: "password",
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		want       string
		wantErr    bool
	}{
		{
			name: "success getting user password",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
						select password from users where id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).
						AddRow("password"))
			},
			want: userSuccess.Password,
		},
		{
			name: "error getting user password - query error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
						select password from users where id = $1
					`).
					WithArgs(1).
					WillReturnError(errors.New("error getting user password"))
			},
			wantErr: true,
		},
		{
			name: "error getting user password - scan error",
			args: args{
				userID: 1,
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectQuery(`
						select password from users where id = $1
					`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"password"}).
						AddRow(nil))
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

			u := Users{
				db: mockDB,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			got, err := u.GetUserPassword(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetUserPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
