#!/bin/bash

old_csv=$(ls | grep csv | tr '\n' ' ')
mv $old_csv archives/$old_csv

cd spider && go run main.go
