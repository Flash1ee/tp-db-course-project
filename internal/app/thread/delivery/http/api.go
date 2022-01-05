package thread_handler

import (
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	thread_usecase "tp-db-project/internal/app/thread/usecase"
	users_handler "tp-db-project/internal/app/users/delivery/http"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	thread_usecase.NotFound: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
var CodeByErrorPost = handler.CodeMap{
	pgx.ErrNoRows:           {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	thread_usecase.NotFound: {http.StatusNotFound, users_handler.NotFound, logrus.WarnLevel},
}
