terraform {
  backend "s3" {
    bucket         = "farhaan-shamsee-ng-voice-backend"
    key            = "eks/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    use_lockfile   = true
  }
}
