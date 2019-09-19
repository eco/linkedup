#!/bin/bash

#cd $TRAVIS_BUILD_DIR
echo "Running linter..."
golangci-lint run
output=$?

if [ $output -eq 0 ]; then
    echo "Linter passed"
    echo "Running ginkgo tests"
    echo "Running tests with mock chain"
    export USE_MOCK_CHAIN=true
    ginkgo ./...
    output=$?
    if [ $output -eq 0 ]; then
        echo "Running tests with LCD chain"
        export USE_MOCK_CHAIN=false
        ginkgo ./...
    else
        echo "Tests on mock chain failed"
        exit $output
    fi

else
    echo "Linter failed "
    exit $output
fi
