package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var ginLambda *ginadapter.GinLambda

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(req)
	if ginLambda == nil {
		ginLambda = ginadapter.New(ginEngine())
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func ginEngine() *gin.Engine {
	app := gin.Default()

	fmt.Println(app)
	app.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	app.POST("/user", func(c *gin.Context) {
		var u user
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, u)
	})

	return app
}

func main() {
	fmt.Println(gin.Mode())
	if gin.Mode() == "release" {
		lambda.Start(lambdaHandler)
	} else {
		app := ginEngine()
		app.Run(":3000")
	}

	// lambda.Start(lambdaHandler)
}
