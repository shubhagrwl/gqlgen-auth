package router

import (
	"context"
	"net/http"
	"time"
	jwt "todo/internal/app/api/middleware"
	"todo/internal/app/api/service"
	"todo/internal/app/service/graph"
	"todo/internal/app/service/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(dependencies service.Services, introspectionEnabled bool) gin.HandlerFunc {
	h := handler.New(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Services: dependencies,
				},
			},
		),
	)
	// h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
	// 	e, ok := err.(error)
	// 	if !ok {
	// 		errString := err.(string)
	// 		return utils.HandleError(ctx, constants.InternalServerError, errors.New(errString))
	// 	}
	// 	return utils.HandleError(ctx, constants.InternalServerError, e)
	// })
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	h.SetQueryCache(lru.New(1000))

	if introspectionEnabled {
		h.Use(extension.Introspection{})
	}
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Init sets up the route for the REST API
func Init(dependencies service.Services) *gin.Engine {
	router := gin.Default()

	// setup cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"Accept", "Accept-CH", "Accept-Charset", "Accept-Datetime", "Accept-Encoding", "Accept-Ext", "Accept-Features", "Accept-Language", "Accept-Params", "Accept-Ranges", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "Access-Control-Expose-Headers", "Access-Control-Max-Age", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Authorization", "Content-Type"}
	corsConfig.AllowAllOrigins = true

	// setup cors middleware
	router.Use(cors.New(corsConfig))

	introspectionEnabled := true

	pingHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
	router.Use(jwt.GinContextToContext())

	router.GET("/playground", playgroundHandler())

	router.GET("/ping", pingHandler)

	router.POST("/query", jwt.NewJWT().Auth(context.Background()), graphqlHandler(dependencies, introspectionEnabled))

	return router
}
