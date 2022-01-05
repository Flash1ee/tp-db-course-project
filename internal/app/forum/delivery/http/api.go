package forum_handler

import (
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	users_handler "tp-db-project/internal/app/users/delivery/http"
	"tp-db-project/internal/app/users/users_usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	pgx.ErrNoRows: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
var CodeByErrorPost = handler.CodeMap{
	pgx.ErrNoRows:          {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	users_usecase.NotFound: {http.StatusNotFound, users_handler.NotFound, logrus.WarnLevel},
}
