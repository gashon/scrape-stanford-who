#!/bin/bash

echo $(ls | grep csv | xargs -I{} cat {} | uniq) > $(date +%Y-%m-%d).csv

old_csv=$(ls | grep csv | tr '\n' ' ')
dir=$(date +%Y-%m-%d)

mv $old_csv $dir/
mv $dir archives/

cd spider && go run main.go

cd ../
git pull
git add .
git commit -m "cron: $(date)"
git push
