package lambda

/*
	This is a project to demonstrate a basic serverless API built
	using Gin for AWS Lambda. It returns fictional corporations
	by default, but it also permits the one administrator to add
	new companies.
*/

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type company struct {
	ID    string `json:"id"`
	Media string `json:"media"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Year  string `json:"year"`
}

var companies = []company{
	{
		ID:    "4",
		Name:  "E Corp",
		Media: "Mr. Robot",
		Type:  "Transnational conglomerate",
		Year:  "2015",
	},
	{
		ID:    "3",
		Name:  "Oceanic Airlines",
		Media: "Lost",
		Type:  "Commercial airline",
		Year:  "2004",
	},
	{
		ID:    "2",
		Name:  "Acme Corporation",
		Media: "Who Framed Rober Rabbit",
		Type:  "Transnational conglomerate",
		Year:  "1988",
	},
	{
		ID:    "1",
		Name:  "Weyland-Yutani Corp",
		Media: "Alien",
		Type:  "Transnational conglomerate",
		Year:  "1979",
	},
}

var adapter *ginadapter.GinLambda

func main() {
	adapter = ginadapter.New(createRouter())

	lambda.Start(handler)

	log.Print("Initialised a Lambda")
}

func handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(context, request)
}

func createRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// Create a router with the default middleware
	router := gin.Default()

	log.Print("Initialised Gin router")

	router.SetTrustedProxies(nil)

	router.GET("/v1/status", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "okay",
		})
	})

	// Initialise the API's methods and routes
	initialiseRoutes(router)

	return router
}

func initialiseRoutes(router *gin.Engine) gin.RouterGroup {
	guest := router.Group("/v1")

	// Add the unauthorised routes
	{
		guest.GET("/companies", func(context *gin.Context) {
			context.IndentedJSON(http.StatusOK, companies)
		})
		guest.GET("/companies/:id", func(context *gin.Context) {
			id := context.Param("id")

			for _, company := range companies {
				if company.ID == id {
					context.IndentedJSON(http.StatusOK, company)
					return
				}
			}

			context.IndentedJSON(http.StatusNotFound, gin.H{"error": "company not found"})
		})
	}

	authorised := router.Group("/v1", gin.BasicAuth(
		gin.Accounts{
			"admin": "password",
		}),
	)

	// Add the authorised routes
	{
		authorised.POST("/companies", func(context *gin.Context) {
			var newCompany company

			if err := context.BindJSON(&newCompany); err != nil {
				return
			}

			companies = append(companies, newCompany)

			context.IndentedJSON(http.StatusCreated, newCompany)
		})
	}

	// Include a tailored help message
	router.GET("/v1/help", func(context *gin.Context) {
		authUser := context.MustGet(gin.AuthUserKey).(string)

		user := context.DefaultQuery(authUser, "guest")

		context.String(http.StatusOK, "Welcome %s.", user)
	})

	return gin.Default().RouterGroup
}

func signalTermination(router *gin.Engine) {
	// Allow graceful termination
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	log.Println("Terminating the server")

	quit <- syscall.SIGTERM

	signal.Notify(quit)

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(context); err != nil {
		log.Fatalf("Waiting; server yet to shutdown: %s", err)
	}

	select {
	case <-context.Done():
		log.Println("Connection terminated gracefully")
	}
}
