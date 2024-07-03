package crmhazar_pkg_http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gitlab.com/salamtm.messenger/slog"
)

func InitSwaggerRoute(router *mux.Router, subRouter string) {
	router.PathPrefix(subRouter + "/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(subRouter+"/swagger/doc.json"), //The swaggerListUrl pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}

func InitSwagger(logger *slog.Logger, args interface{}) {

	//List of Swagger file paths
	swaggerFilePaths := []string{"docs.go", "swagger.json", "swagger.yaml"}

	// Iterate through each Swagger file
	for _, filePath := range swaggerFilePaths {
		filePath = "./docs/" + filePath
		// Read the content of the Swagger file
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Error("Error reading file %s: %v\n", filePath, err)
			continue
		}

		newContent := changeArgs(content, args)

		// Write the modified content back to the Swagger file
		err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			logger.Error("Error writing file %s: %v\n", filePath, err)
			continue
		}

		logger.Info("Updated Swagger file %s\n", filePath)
	}

	logger.Info("All Swagger files updated successfully!")
}

func changeArgs(content []byte, args interface{}) string {

	newContent := string(content)

	val := reflect.ValueOf(args)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()

		newContent = strings.ReplaceAll(newContent, "{{ ."+field.Name+" }}", fmt.Sprintf("%v", fieldValue))
	}

	return newContent
}
