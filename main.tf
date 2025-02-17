resource "aws_key_pair" "ssh_key" {
  key_name   = "mypc-ssh-key"
  public_key = file("~/.ssh/id_rsa.pub") 
}

resource "aws_security_group" "allow_ssh" {
  name        = "allow_ssh"
  description = "Allow SSH inbound traffic"

  ingress {
    description = "SSH from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] 
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "coderunner_single" {
  ami           = "ami-032a56ad5e480189c"
  instance_type = "t3.micro"
  key_name      = aws_key_pair.ssh_key.key_name
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]
  associate_public_ip_address = true  

  tags = {
    Name = "CoderunnerMainInstance"
  }
}

# Output the public IP for easy access
output "public_ip" {
  value = aws_instance.coderunner_single.public_ip
}
