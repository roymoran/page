package providers

// Terraform definition required to deployed
// necessary resources to host a static site on
// AWS. Includes creation of S3, S3 Bucket, CloudFront,
// and CloudFront distribution
var AwsTerraformDefinition string = `terraform {
	required_providers {
	  aws = {
		source  = "hashicorp/aws"
		version = "3.25.0"
	  }
	}
  }
  
  provider "aws" {
	profile = "default"
	region  = "us-east-2"
  }

  locals {
	s3_origin_id = "myS3Origin"
	s3_bucket_name = "pagecli-2827005998"
  }
  
  resource "aws_s3_bucket" "b" {
	bucket = local.s3_bucket_name
  
	website {
	  index_document = "index.html"
	  error_document = "index.html"
	}
  }
  
  resource "aws_s3_bucket_object" "s3website" {
	bucket       = local.s3_bucket_name
	key          = "mywebsite9.com/index.html"
	source       = "/Users/roymoran/.pagecli/deploy/index.html"
	acl          = "public-read"
	content_type = "text/html"
  
	# The filemd5() function is available in Terraform 0.11.12 and later
	# For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
	# etag = "${md5(file("path/to/file"))}"
	etag       = filemd5("/Users/roymoran/.pagecli/deploy/index.html")
	depends_on = [aws_s3_bucket.b]
  }
  
  resource "aws_cloudfront_distribution" "s3_distribution" {
	origin {
	  domain_name = aws_s3_bucket.b.bucket_regional_domain_name
	  origin_path = "/mywebsite9.com"
	  origin_id   = local.s3_origin_id
	}
  
	enabled             = true
	is_ipv6_enabled     = true
	default_root_object = "index.html"
  
	# TODO: ADD CNAMES
	# aliases = ["mysite.example.com", "yoursite.example.com"]
  
	default_cache_behavior {
	  allowed_methods  = ["GET", "HEAD"]
	  cached_methods   = ["GET", "HEAD"]
	  target_origin_id = local.s3_origin_id

	  forwarded_values {
		query_string = true
  
		cookies {
		  forward = "none"
		}
	  }

	  viewer_protocol_policy = "redirect-to-https"
	  min_ttl                = 0
	  default_ttl            = 3600
	  max_ttl                = 86400
	}
  
	price_class = "PriceClass_All"

	restrictions {
		geo_restriction {
		  restriction_type = "none"
		}
	}

	viewer_certificate {
	  cloudfront_default_certificate = true
	}
  
	depends_on = [aws_s3_bucket.b]
  }
`

type TerraformTemplate struct {
	Terraform RequiredProviders         `json:"terraform"`
	Provider  map[string]ProviderConfig `json:"provider"`
}

type RequiredProviders struct {
	RequiredProvider map[string]Provider `json:"required_providers"`
}

type Provider struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

type ProviderConfig struct {
	Profile string `json:"profile"`
	Region  string `json:"region"`
}
