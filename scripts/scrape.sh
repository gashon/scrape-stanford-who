#!/bin/bash

old_csv=$(ls | grep csv | tr '\n' ' ')
mv $old_csv archives/

cd spider && go run main.go

<<<<<<< HEAD
cd ../
=======
git pull
>>>>>>> cb3e661 (fix: git pull)
git add .
git commit -m "cron: $(date)"
git push
