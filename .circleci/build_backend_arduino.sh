#!/usr/bin/env bash
mkdir backend_arduino_build
cmake -B backend_arduino_build -H cmd/backsrv/backend_arduino/
make -C backend_arduino_build -j4
