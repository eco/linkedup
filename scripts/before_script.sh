#!/bin/bash


#install test framework
echo "Installing ginko and gomega"
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...

requiredVersion="v1.18.0"
needsInstall=0
#installing the linter
if ! [ -x "$(command -v golangci-lint)" ]; then
  needsInstall=1
else
  ver=`golangci-lint --version`

  if [[ $ver != *"${requiredVersion}"* ]]; then
    needsInstall=1
  fi
fi

if [ $needsInstall -eq 1 ]; then
  echo "Installing golangci-lint"
  go get github.com/golangci/golangci-lint/cmd/golangci-lint@${requiredVersion}
  echo "Finished installing gometaliner"
else
  echo "golangci-lint already installed"
  which golangci-lint
fi


# Stoyans git api token saved as env variable in travis
echo "Adding git global for the go mod tool to access private repos"
git config --global url."https://$GITHUB_KEY:@github.com/".insteadOf "https://github.com/"
