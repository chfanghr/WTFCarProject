#!/usr/bin/env bash
mkdir backend_arduino_build
cd backend_arduino_build;cmake cmd/backsrv/backend_arduino/
cd backend_arduino_build;make -C backend_arduino_build -j4
