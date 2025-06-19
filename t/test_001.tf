resource aws_instance "web" {
  ami = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}

module "vpc" {
  source = "../t/module/test"
}

module "vpc2" {
  source = "git@github.com/terraform-aws-modules/terraform-aws-vpc.git?ref=v5.21.1"
}
