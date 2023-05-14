locals {
  source_path               = "${path.module}/../cmd"
  binary_path               = "${path.module}/out/main"
  build_binary_command      = "env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${local.binary_path} ${local.source_path}"
  zip_path                  = "${path.module}/out/main.zip"
  aws_resource_name_postfix = "Terraform"
}

# Build + ZIP file =================================================================================

resource "null_resource" "create_binary" {
  count = var.recreate_zip_file ? 1 : 0
  triggers = {
    condition = timestamp()
  }

  provisioner "local-exec" {
    command = local.build_binary_command
  }
}

data "archive_file" "lambda_zip" {
  depends_on  = [null_resource.create_binary]
  type        = "zip"
  source_file = local.binary_path
  output_path = local.zip_path
}

# Database =========================================================================================

resource "aws_dynamodb_table" "monthly_planning" {
  name           = "MonthlyPlanning${local.aws_resource_name_postfix}"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "slot"
  attribute {
    name = "slot"
    type = "S"
  }
}

# Permissions ======================================================================================

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda" {
  name               = "MonthlyPlanningLambdaRole${local.aws_resource_name_postfix}"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

data "aws_iam_policy_document" "table_access" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = [aws_dynamodb_table.monthly_planning.arn]
  }
}

resource "aws_iam_role_policy" "lambda_dynamodb" {
  name   = "MonthlyPlanningDynamoDB${local.aws_resource_name_postfix}"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.table_access.json
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  count      = var.enable_cloudwatch_logs ? 1 : 0
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Lambda function + URL ============================================================================

resource "aws_lambda_function" "monthly_planning" {
  depends_on       = [data.archive_file.lambda_zip]
  filename         = local.zip_path
  function_name    = "MonthlyPlanning${local.aws_resource_name_postfix}"
  role             = aws_iam_role.lambda.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  environment {
    variables = {
      AUTH_USERNAME = var.auth_username
      AUTH_PASSWORD = var.auth_password
      TABLE_NAME    = aws_dynamodb_table.monthly_planning.name
    }
  }
}

resource "aws_lambda_function_url" "monthly_planning" {
  function_name      = aws_lambda_function.monthly_planning.function_name
  authorization_type = "NONE"

  cors {
    allow_origins = var.cors_allow_origins
    allow_methods = ["*"]
    allow_headers = ["*"]
  }
}
