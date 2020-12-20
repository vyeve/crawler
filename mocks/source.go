package mocks

//go:generate mockgen -destination=./io-writer.go -package=mocks io Writer
//go:generate mockgen -destination=./read-closer.go -package=mocks io ReadCloser
