#!/bin/bash

# https://github.com/awsdocs/aws-lambda-developer-guide
# 
# Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy of this
# software and associated documentation files (the "Software"), to deal in the Software
# without restriction, including without limitation the rights to use, copy, modify,
# merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
# permit persons to whom the Software is furnished to do so.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
# INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
# PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
# OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
# SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

set -eo pipefail
STACK=monthly-planning-backend
ARTIFACT_BUCKET=$(cat bucket-name.txt)
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o out/main ./cmd
aws cloudformation package --template-file out/template.yml --s3-bucket $ARTIFACT_BUCKET --output-template-file out/out.yml
aws cloudformation deploy --template-file out/out.yml --stack-name $STACK --capabilities CAPABILITY_NAMED_IAM
