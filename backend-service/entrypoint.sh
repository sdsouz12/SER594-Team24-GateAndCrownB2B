#!/bin/sh
set -e
echo "Running database migrations..."
./migrate up
echo "Starting backend server..."
exec ./server
