#!/bin/bash
for ((i=1; i <= $2 ; i++))
do
  $GOPATH/bin/cmd $GOPATH/src/github.com/PuerkitoBio/specter/cmd/examples/$1.vm
done
