#!/bin/bash
git pull

old_csv=$(ls | grep csv | tr '\n' ' ')
dir=$(date -d "yesterday" +%Y-%m-%d)

mkdir $dir
mv $old_csv $dir/
mv $dir archives/

cd spider && go run main.go

cd ../
c=$(date +%Y-%m-%d).csv
echo $(ls | grep csv | xargs -I{} cat {} | uniq) > $c

git add $(ls | grep csv | tr '\n' ' ') archives/$dir 
git commit -m "cron: $(date)"
git push

rm -rf archives/$dir