# se crea el buket 
resource "aws_s3_bucket" "frontend" {
  bucket = "frontend-${random_id.suffix.hex}"
}

# bloqueo politicas de acceso
resource "aws_s3_bucket_public_access_block" "block_public" {
  bucket = aws_s3_bucket.frontend.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# se suben los archivos de distr
resource "aws_s3_object" "frontend_files" {
  for_each     = fileset("../../Frontend/vue-frontend/dist", "**/*")
  bucket       = aws_s3_bucket.frontend.id
  key          = each.value
  source       = "../../Frontend/vue-frontend/dist/${each.value}"
  content_type = lookup({ 
    "html" = "text/html",
    "css"  = "text/css",
    "js"   = "application/javascript",
    "ico"  = "image/x-icon",
    "png"  = "image/png"
  }, split(".", each.value)[1], "application/octet-stream")
}