#!/bin/bash



cd "$(dirname "$0")"

mkdir -p build

rm -rf build

cd backend && go build -o ../build/network_monitor_backend
cd ..

cd frontend/local_network_monitor

flutter build linux --release

cp -r build/linux/x64/release/bundle/* ../../build/
