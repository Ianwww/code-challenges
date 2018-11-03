# Create keypair
Use your existing pub/priv keys or create new ones using the command below.  Add public key to field "public_key" of the "aws_key_pair" section of nginx-setup.tf file
```
ssh-keygen -f mykey
```

#  Using the aws console IAM, create a user account with programatic access and policy "AmazonEC2FullAccess." Add your aws account keys to .tf file provider section. 
```
access_key = <enter key>
secret_key = <enter key>
```

# Build cluster
* Install terraform locally
* Download "ian" folder with terraform files
* From folder where .tf file is, run
  - terraform init
  - terraform plan 
  - terraform apply 
  - type yes when prompted
    _it may take up to 15 minutes after command returns for the cluster to be ready for you to view the solution_
    - get public ip address of cluster from the terraform command output and paste it into your browser.
    _example : http://ec2-54-158-28-233.compute-1.amazonaws.com_

