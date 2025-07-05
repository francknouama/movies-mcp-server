#!/bin/bash

# Build script for BDD tests
set -e

echo "Building BDD test structure..."

# Check if we can build the context
echo "Testing context package..."
go build ./context

# Check if we can build the support package  
echo "Testing support package..."
go build ./support

# Check if we can build the steps package
echo "Testing steps package..."
go build ./steps

echo "All BDD packages compiled successfully!"

# Try to compile the test
echo "Testing BDD test compilation..."
go test -c .

echo "BDD structure compilation successful!"