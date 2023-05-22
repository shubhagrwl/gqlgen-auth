package context

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"todo/internal/app/service/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey is an internal type used for holding values in a `context.Context`
type ContextKey string

const (
	keyGinContextKey ContextKey = "auth-gin-context"
	keyUserID        ContextKey = "user-id"
	keyRequest       ContextKey = "auth-req"
	keyRawRequest    ContextKey = "auth-raw-req"
)

// CreateContextFromGinContext creates the context object from gin context
func CreateContextFromGinContext(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), keyGinContextKey, c)
}

// WithUserID is used to configure a context to add user id
func WithUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(contextOrBackground(parent), keyUserID, userID)
}

// UserIDFromContext returns the user id from the context
func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	log := logger.Logger(ctx)
	fmt.Println(ctx.Value(keyUserID))
	userID, ok := ctx.Value(keyUserID).(string)
	if !ok {
		log.Error("No user id present in context...")
		return uuid.Nil, errors.New("no id present in context")
	}

	id := uuid.MustParse(userID)
	return id, nil
}

// GinContextFromContext returns the gin context from context
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(keyGinContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// WithRequest adds the user data to context
func WithRequest(ctx context.Context, c *gin.Context) context.Context {
	body, _ := ioutil.ReadAll(c.Request.Body)
	ctx = context.WithValue(ctx, keyRequest, string(body))
	ctx = context.WithValue(ctx, keyRawRequest, *c.Request)
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	return ctx

}

// RequestFromContext returns the user token from the context
func RequestFromContext(ctx context.Context) (string, error) {
	log := logger.Logger(ctx)

	req, ok := ctx.Value(keyRequest).(string)
	if !ok {
		log.Error("No request present in context...")
		return "", errors.New("no request present in context")

	}
	return req, nil

}

// RawRequestFromContext returns the user token from the context
func RawRequestFromContext(ctx context.Context) (*http.Request, error) {
	log := logger.Logger(ctx)

	req, ok := ctx.Value(keyRawRequest).(http.Request)
	if !ok {
		log.Error("No request present in context...")
		return nil, errors.New("no request present in context")

	}
	return &req, nil

}

// contextOrBackground returns the given context if it is not nil.
// Returns context.Background() otherwise.
func contextOrBackground(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return context.Background()
}
