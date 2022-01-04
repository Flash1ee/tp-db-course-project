package post_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	post_repository "tp-db-project/internal/app/post/repository"
	"tp-db-project/internal/pkg/handler"
)

var CodeByErrorGet = handler.CodeMap{
	post_repository.NotFound: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
var CodeByErrorPost = handler.CodeMap{
	post_repository.NotFound: {http.StatusNotFound, ForumNotFound, logrus.WarnLevel},
}
