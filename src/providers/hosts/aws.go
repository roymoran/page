package providers

import (
	"fmt"
)

type AmazonWebServices struct {
	Infrastructure string
}

var AwsTerraformProvider = `
terraform {
	required_providers {
		aws = {
			source = "hashicorp/aws"
			version = "3.25.0"
		}
	}
}
`

var AmazonWebServicesConfig = AmazonWebServices{
	Infrastructure: `
	terraform {
		required_providers {
		  aws = {
			source = "hashicorp/aws"
			version = "3.25.0"
		  }
		}
	  }

	  provider "aws" {
		profile = "default"
		region  = "us-east-2"
	  }

	  resource "aws_s3_bucket" "b" {
		bucket = "my-tf-test-bucket"
		acl    = "private"
	  
		tags = {
		  Name        = "My bucket"
		  Environment = "Dev"
		}
	  }
	`,
}

func (aws AmazonWebServices) Deploy() bool {
	return true
}

func (aws AmazonWebServices) ConfigureHost() bool {
	fmt.Println("configured aws host")
	return true
}

func (aws AmazonWebServices) HostProviderDefinition() string {
	return AwsTerraformProvider
}
