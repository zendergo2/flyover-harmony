#!/bin/bash

rm build/*
go build -o build/
./build/flyover-harmony
