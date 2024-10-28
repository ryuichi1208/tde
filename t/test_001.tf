resource aws_instance "web" {
  ami = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}

module "vpc" {
  source = "./vpc"
}
