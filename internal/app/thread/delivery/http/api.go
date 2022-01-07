package thread_handler

import (
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	thread_usecase "tp-db-project/internal/app/thread/usecase"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	thread_usecase.NotFound: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
var CodeByErrorPost = handler.CodeMap{
	thread_usecase.CreatedEmpty:   {http.StatusCreated, thread_usecase.CreatedEmpty, logrus.InfoLevel},
	thread_usecase.AuthorNotFound: {http.StatusNotFound, thread_usecase.AuthorNotFound, logrus.ErrorLevel},
	thread_usecase.ParentInvalid:  {http.StatusConflict, thread_usecase.ParentInvalid, logrus.WarnLevel},
	thread_usecase.ParentNotFound: {http.StatusConflict, thread_usecase.ParentNotFound, logrus.WarnLevel},
	pgx.ErrNoRows:                 {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
	thread_usecase.NotFound:       {http.StatusNotFound, thread_usecase.NotFound, logrus.WarnLevel},
}
