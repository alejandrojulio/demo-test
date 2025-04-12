
output "bucket_name" {
  value = aws_s3_bucket.frontend.bucket
}

# outputs.tf
output "cloudfront_url" {
  value = aws_cloudfront_distribution.s3_distribution.domain_name
}
