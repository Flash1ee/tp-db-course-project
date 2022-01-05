package users_handler

import (
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/internal/app/users/users_usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	pgx.ErrNoRows: {http.StatusNotFound, NotFound, logrus.InfoLevel},
}
var CodeByErrorPost = handler.CodeMap{
	pgx.ErrNoRows:                  {http.StatusNotFound, NotFound, logrus.InfoLevel},
	users_usecase.AlreadyExistsErr: {http.StatusConflict, ConflictErr, logrus.WarnLevel},
}
