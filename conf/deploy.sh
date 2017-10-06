#!/usr/bin/env bash

echo "Building..."
go-bindata ./static/...
go build -ldflags "-X main.version=$(git rev-parse HEAD)"

echo "Stopping app"
ssh -p 2280 web@storage.bad.mn sudo systemctl stop pagaraprata

echo "Uploading..."
scp -P 2280 pagaraprata pagaraprata@storage.bad.mn:/home/pagaraprata/app/pagaraprata

echo "Starting"
ssh -p 2280 web@storage.bad.mn sudo systemctl daemon-reload
ssh -p 2280 web@storage.bad.mn sudo systemctl start pagaraprata
