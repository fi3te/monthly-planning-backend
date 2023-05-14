variable "recreate_zip_file" {
  type        = bool
  default     = true
  description = "Defines whether the zip file with an executable binary file for the Lambda function should be recreated."
}

variable "enable_cloudwatch_logs" {
  type        = bool
  default     = false
  description = "Defines whether the lambda function should create and write to a CloudWatch log group."
}

variable "auth_username" {
  type        = string
  default     = "admin"
  description = "HTTP Basic username"
}

variable "auth_password" {
  type        = string
  default     = "password"
  description = "HTTP Basic password"
}

variable "cors_allow_origins" {
  type    = list(any)
  default = ["*"]
}
