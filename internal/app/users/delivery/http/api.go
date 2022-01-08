package users_handler

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/internal/app/users/users_usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	users_usecase.NotFound: {http.StatusNotFound, NotFound, logrus.InfoLevel},
	pgx.ErrNoRows:          {http.StatusNotFound, NotFound, logrus.InfoLevel},
}
var CodeByErrorPost = handler.CodeMap{
	users_usecase.NotFound:         {http.StatusNotFound, NotFound, logrus.InfoLevel},
	users_usecase.ConstraintError:  {http.StatusConflict, ConflictErr, logrus.InfoLevel},
	users_usecase.AlreadyExistsErr: {http.StatusConflict, ConflictErr, logrus.WarnLevel},
}
