package Middlewares

import (
	"Kaban/internal/DomainLevel"
	"fmt"
	"log/slog"
	"net/http"
)

func CheckPostRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			slog.Error("The request with a invalid method", slog.Group("Request detail"),
				slog.String("Method", request.Method), slog.String("The URL", request.RequestURI))
			writer.Header().Set(DomainLevel.ContentType, DomainLevel.Json)
			writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(writer, "The method isn't right")
			return
		}
		next.ServeHTTP(writer, request)
	})
}
