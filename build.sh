#!/bin/bash

go-bindata fonts/*;
go build .;
./go-svgpng