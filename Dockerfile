FROM golang:latest AS build-stage
LABEL authors="pasabaranov"
WORKDIR /bu
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod tidy
COPY . .
COPY ./cmds  ./
COPY ./DipendsInjective ./
COPY ./DomainLevel ./
COPY ./Dto ./
COPY ./InfrastructureLevel ./

RUN echo "We finished copying"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./main.go


FROM alpine:latest AS final-stage
WORKDIR /app
COPY --from=build-stage /bu/app ./
COPY --from=build-stage /bu/cmds ./
COPY --from=build-stage /bu/DipendsInjective ./
COPY --from=build-stage /bu/Dto ./
COPY --from=build-stage /bu/InfrastructureLevel ./
COPY --from=build-stage /bu/DomainLevel ./
EXPOSE 80
RUN echo "We completed"
CMD ["./app"]