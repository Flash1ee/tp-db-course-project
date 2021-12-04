package users_repository

import (
	"database/sql"
	"tp-db-project/internal/app/users/models"
)

const (
	queryCreateUser = "INSERT INTO users(nickname, fullname, about, email) VALUES($1, $2, $3, $4)"
	queryGetUser    = "SELECT nickname, fullname, about, email FROM users where nickname = $1"
	queryUpdateUser = "UPDATE users SET " +
		"fullname = COALESCE(NULLIF(TRIM($1), ''), fullname)," +
		"about = COALESCE(NULLIF(TRIM($2), ''), about), " +
		"email = COALESCE(NULLIF(TRIM($3), ''), email) " +
		"where nickname = $4" +
		"RETURNING fullname, about, email;"
)

type UsersRepository struct {
	conn *sql.DB
}

func NewUsersRepository(conn *sql.DB) *UsersRepository {
	return &UsersRepository{
		conn: conn,
	}
}

func (r *UsersRepository) Create(user *models.User) error {
	_, err := r.conn.Exec(queryCreateUser, user.Nickname, user.FullName, user.About, user.Email)

	return err
}

func (r *UsersRepository) Get(nickname string) (*models.User, error) {
	user := &models.User{}
	err := r.conn.QueryRow(queryGetUser, nickname).Scan(user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UsersRepository) Update(user *models.User) (*models.User, error) {
	newUser := &models.User{}
	if err := r.conn.QueryRow(queryUpdateUser, user.FullName, user.About, user.Email, user.Nickname).
		Scan(&newUser.FullName, &newUser.About, &newUser.Email); err != nil {
		return nil, err
	}

	return newUser, nil
}
