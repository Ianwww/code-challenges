#!/bin/sh

sudo yum install epel-release -y
#sudo yum update -y
sudo yum install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
#sudo yum clean all
NGINX_HOSTNAME=`hostname -s`
echo "$NGINX_HOSTNAME"
echo "<html></body>ByteCubed Challenge $NGINX_HOSTNAME </body></html>" > /usr/share/nginx/html/index.html

#firewall-cmd --zone=public --permanent --add-service=http
#firewall-cmd --zone=public --permanent --add-service=https
#firewall-cmd --reload