#!/bin/bash

old_csv=$(ls | grep csv | tr '\n' ' ')
mv $old_csv archives/

cd spider && go run main.go

cd ../
git pull
git add .
git commit -m "cron: $(date)"
git push
