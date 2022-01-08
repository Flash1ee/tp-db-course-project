package forum_handler

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"net/http"
	forum_usecase "tp-db-project/internal/app/forum/usecase"
	users_handler "tp-db-project/internal/app/users/delivery/http"
	"tp-db-project/internal/app/users/users_usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	forum_usecase.ForumNotFound: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	pgx.ErrNoRows:               {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
var CodeByErrorPost = handler.CodeMap{
	forum_usecase.ConstraintError: {http.StatusConflict, forum_usecase.AlreadyExists, logrus.WarnLevel},
	forum_usecase.AlreadyExists:   {http.StatusConflict, forum_usecase.AlreadyExists, logrus.WarnLevel},
	forum_usecase.ForumNotFound:   {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	pgx.ErrNoRows:                 {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	users_usecase.NotFound:        {http.StatusNotFound, users_handler.NotFound, logrus.WarnLevel},
}
