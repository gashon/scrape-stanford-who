#!/bin/bash
git pull

old_csv=$(ls | grep csv | tr '\n' ' ')
dir=$(date -v -1d '+%Y-%m-%d')

mkdir $dir
mv $old_csv $dir/
mv $dir archives/

cd spider && go run main.go

cd ../
echo $(ls | grep csv | xargs -I{} cat {} | uniq) > $(date +%Y-%m-%d).csv

git add .
git commit -m "cron: $(date)"
git push
