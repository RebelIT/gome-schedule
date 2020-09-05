#!/bin/bash

echo "Getting Dependencies"
go get

echo "Building the program"
go build -o main .

echo "Running the program"
exec ./main $PARAMS $@