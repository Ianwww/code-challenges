###setup provider
provider "aws" {
  access_key = <redacted>
  secret_key = <redacted>
  region     = "us-east-1"
}

###networking
#VPC
resource "aws_vpc" "my_vpc" {
  cidr_block = "172.16.0.0/16"
  enable_dns_hostnames = true
  tags {
    Name = "bytecubed-ws-example"
  }
}

#Subnet in VPC, webserver to be deployed in this subnet
resource "aws_subnet" "public_subnet" {
  vpc_id = "${aws_vpc.my_vpc.id}"
  cidr_block = "172.16.50.0/24"
  availability_zone = "us-east-1a"
  tags {
    Name = "bytecubed-ws-example"
  }
}

#gateway for the vpc
resource "aws_internet_gateway" "bytecubed_gw" {
  vpc_id = "${aws_vpc.my_vpc.id}"

  tags {
    Name = "VPC IGW"
  }
}

#public facing elastic ip
resource "aws_eip" "lb" {
  instance = "${aws_instance.bytecubed_instance.id}"
  vpc      = true
}

# network interface for ec2 instance
resource "aws_network_interface" "bytecubed_net" {
  subnet_id = "${aws_subnet.public_subnet.id}"
  private_ips = ["172.16.50.100"]
  tags {
    Name = "primary_network_interface"
  }
}

# Define the route table
resource "aws_route_table" "vpc_rt" {
  vpc_id = "${aws_vpc.my_vpc.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.bytecubed_gw.id}"
  }

  tags {
    Name = "vpc rt"
  }
}

# Assign the route table to the public Subnet
resource "aws_route_table_association" "web-public-rt" {
  subnet_id = "${aws_subnet.public_subnet.id}"
  route_table_id = "${aws_route_table.vpc_rt.id}"
}

# Define the security group for public subnet
resource "aws_security_group" "sg_nginx" {
  name = "vpc_nginx"
  description = "Allow incoming HTTP connections to the vpc public subnet"

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks =  ["0.0.0.0/0"]
  }
  vpc_id="${aws_vpc.my_vpc.id}"

  tags {
    Name = "nginx sg"
  }
}

###Webserver setup
#create a ec2 instance and install nginx
resource "aws_instance" "bytecubed_instance" {
  ami           = "ami-4bf3d731"
  instance_type = "t2.micro"
  subnet_id = "${aws_subnet.public_subnet.id}"
  vpc_security_group_ids = ["${aws_security_group.sg_nginx.id}"]
  associate_public_ip_address = true
  source_dest_check = false
  user_data = "${file("./install.sh")}"
  key_name = "bytecubed-ssh"

  tags {
    Name = "bytecubed-nginx"
  }
}

resource "aws_key_pair" "default" {
  key_name = "bytecubed-ssh"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDdI1OQ1b1Hmwg8NRMka/WcvIQ/oxRDmInEG5RkQZCoZUUSpIU96pG30OaH0sX1jSMK2frSeTGTlsp1uamUw3qDYqbN3xKZ0+aJvc7xu/R6kH+QeBPT6I+uVYe8Wvx+rqv2cmbxShf5256PYnxdXl1sm3AW8FmoVHIMzfP9UDEFd9bBMBbPEDShQj/vnNLpiXfGk/Gc/MCmvQEMq2wEyY4Yk+fOF+5ZhdE64tUY1q6Bizix4GLIwgzbXskWClBSNXSEiVcTuYnLllkAc4X/oQVuXDPE+8E7D24E1nE3zSdqWC1+CcuyhObzXVPIGMIY/9fvlsQ9Ob7//2+50LG5vWki2zFvZ8ln30nVAgW4U15nxZUP9kzidGkU2R560dHAYeeZGh4GhZWuFan8U++8TgeAvpb5i8KQTQSOC5aCMr1Zlch4Y7wK1Xg1OYvy9ea72qNYKG4xEsnyn63TAueP580rVUHkBUsyeAwEuITZKeEAar0P965aKVqaJg4/5vnXT0rsh3Iqe24kes4VrptDbfVTUWJYHAT7yzpnndSLC5+aRAIjLAPOE8Fst+SjOldXzwPAvygqplhQ4nw+EDW2M53ah7FCZnaoydlApEgLUC1O1F+SjzexNILFmXkgkdaymlp9pfyTNoQhunUSOtENzaCXyhLMbXGP+wxAGkaxfOCg9w== ianbw2@gmail.com"
}

###alternative do in ecs cluster and deploy with docker