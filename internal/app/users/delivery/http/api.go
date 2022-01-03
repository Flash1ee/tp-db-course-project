package users_handler

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"net/http"
	"tp-db-project/internal/app/users/users_usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	sql.ErrNoRows: {http.StatusNotFound, NotFound, logrus.InfoLevel},
}
var CodeByErrorPost = handler.CodeMap{
	sql.ErrNoRows:                  {http.StatusNotFound, NotFound, logrus.InfoLevel},
	users_usecase.AlreadyExistsErr: {http.StatusConflict, ConflictErr, logrus.WarnLevel},
}
