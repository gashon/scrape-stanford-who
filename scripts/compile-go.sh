#!/bin/bash

# Define the path to the main Go file
MAIN_FILE="./spider/main.go"

# Check if the main file exists
if [ ! -f "$MAIN_FILE" ]; then
  echo "Error: $MAIN_FILE not found"
  exit 1
fi

# Check if the project is already compiled
if [ -f "$MAIN_FILE.out" ]; then
  echo "Project already compiled"
else
  echo "Compiling project..."
  go build "$MAIN_FILE"
  if [ $? -eq 0 ]; then
    echo "Compilation successful"
  else
    echo "Error: compilation failed"
    exit 1
  fi
fi