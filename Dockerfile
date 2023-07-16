FROM node:16.15-alpine
WORKDIR /app
COPY ./front-end ./
RUN npm run i
RUN npm run build

FROM golang:1.19.4-alpine
WORKDIR /app
COPY go.mod ./
COPY config.json ./
RUN go mod download
ADD . .
COPY --from=0 /app/build ./front-end/build
RUN go build -o /server-status-monitoring

EXPOSE 3010

CMD [ "/server-status-monitoring" ]