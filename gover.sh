#!/bin/bash

dirs=$(find . -type d | grep -v .git | grep ./)

for dir in $dirs; do
  go test -coverprofile=${dir:2}_tmp.coverprofile ${dir}
done

${GOPATH}/bin/gover
rm *_tmp.coverprofile
