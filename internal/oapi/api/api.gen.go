// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all tasks
	// (GET /tasks)
	GetTasks(c *gin.Context)
	// Create task
	// (POST /tasks)
	PostTasks(c *gin.Context)
	// Delete a task by ID
	// (DELETE /tasks/{id})
	DeleteTasksId(c *gin.Context, id uint64)
	// Update a task by ID
	// (PATCH /tasks/{id})
	PatchTasksId(c *gin.Context, id uint64)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetTasks operation middleware
func (siw *ServerInterfaceWrapper) GetTasks(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetTasks(c)
}

// PostTasks operation middleware
func (siw *ServerInterfaceWrapper) PostTasks(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostTasks(c)
}

// DeleteTasksId operation middleware
func (siw *ServerInterfaceWrapper) DeleteTasksId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id uint64

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteTasksId(c, id)
}

// PatchTasksId operation middleware
func (siw *ServerInterfaceWrapper) PatchTasksId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id uint64

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PatchTasksId(c, id)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/tasks", wrapper.GetTasks)
	router.POST(options.BaseURL+"/tasks", wrapper.PostTasks)
	router.DELETE(options.BaseURL+"/tasks/:id", wrapper.DeleteTasksId)
	router.PATCH(options.BaseURL+"/tasks/:id", wrapper.PatchTasksId)
}
