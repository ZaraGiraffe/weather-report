// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"example.com/weather-report/restapi/operations"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"

	"example.com/weather-report/storage"
	"example.com/weather-report/handlers"
	"example.com/weather-report/config"
)

//go:generate swagger generate server --target ..\..\weather-report --name WeatherForecastAPI --spec ..\task.yaml --principal interface{}

func configureFlags(api *operations.WeatherForecastAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.WeatherForecastAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.UrlformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// subscription.SubscribeMaxParseMemory = 32 << 20

	storage := storage.NewStorageConnection()
	conf := config.GetConfig()

	api.SubscriptionConfirmSubscriptionHandler = subscription.ConfirmSubscriptionHandlerFunc(handlers.ConfirmSubscriprionHandlerWrapper(storage))
	api.WeatherGetWeatherHandler = weather.GetWeatherHandlerFunc(handlers.GetWeatherHandlerWrapper(&conf.WeatherApiConfig))
	api.SubscriptionSubscribeHandler = subscription.SubscribeHandlerFunc(handlers.SubscribeHandlerWrapper(conf, storage))
	api.SubscriptionUnsubscribeHandler = subscription.UnsubscribeHandlerFunc(handlers.UnsubscribeHandlerWrapper(storage))

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
