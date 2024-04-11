FROM golang:1.21
LABEL authors="Pasca11"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

EXPOSE 3000

CMD make run