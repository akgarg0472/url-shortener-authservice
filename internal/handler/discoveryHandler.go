package handler

import (
	"net/http"
	"strconv"

	build "github.com/akgarg0472/urlshortener-auth-service/build"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

func DiscoveryInfoHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	port, _ := strconv.Atoi(utils.GetEnvVariable("SERVER_PORT", "8081"))

	infoResponse := map[string]interface{}{
		"build": map[string]interface{}{
			"buildTime":     build.BuildTime,
			"go.version":    build.GoVersion,
			"compiler.os":   build.OS,
			"compiler.arch": build.Arch,
		},
		"app": map[string]interface{}{
			"version":     build.AppVersion,
			"artifact":    "AuthService",
			"name":        "Auth Service",
			"description": "Auth Service for URL Shortener project",
		},
		"runtime": map[string]interface{}{
			"go": map[string]interface{}{
				"version": build.GoVersion,
				"arch":    build.Arch,
			},
			"port": port,
			"ip":   utils.GetHostIP(),
		},
	}

	sendResponseToClient(responseWriter, "", infoResponse, nil, 200)
}

func DiscoveryHealthHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	sendResponseToClient(responseWriter, "", map[string]interface{}{"status": "UP"}, nil, 200)
}
