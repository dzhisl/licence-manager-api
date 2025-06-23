package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FormResponse(msg string) (int, gin.H) {
	return 200, gin.H{"message": msg}
}

func FormErrResponse(status int, err string) (int, gin.H) {
	return status, gin.H{"error": err}
}

func FormInvalidRequestResponse() (int, gin.H) {
	return FormErrResponse(http.StatusBadRequest, "invalid request")
}

func FormInternalErrResponse() (int, gin.H) {
	return FormErrResponse(http.StatusInternalServerError, "internal server error")
}
