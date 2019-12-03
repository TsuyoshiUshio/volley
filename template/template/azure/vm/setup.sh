# ********************************************************************
# * Script run by the Custom Script Extension on the provisioning VM *
# ********************************************************************
# Set Azure Credentials by reading the command line arguments

echo "############### Adding package respositories ###############"
# Get the Docker GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg 2>&1 | sudo apt-key add -

# Add Docker source
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

echo "############### Installing Packages ###############"

sudo DEBIAN_FRONTEND=noninteractive apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y curl
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y gnupg-agent
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y software-properties-common
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y docker-ce docker-ce-cli containerd.io

#Add user to docker usergroup
sudo DEBIAN_FRONTEND=noninteractive apt-get remove -y unscd
sudo usermod -aG docker azureuser

#Holding walinuxagent before upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-mark hold walinuxagent
sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y

# log directory creation
mkdir -p /home/azureuser/logs

# /bin/bash -c 'docker run -d --name docker-daemon --privileged docker:stable-dind &'
# /bin/bash -c 'docker run -v /home/nginx/config:/home/nginx/config -v /home/nginx/contents:/home/nginx/contents -v /home/azureuser/logs:/home/azureuser/logs -v /var/run/docker.sock:/var/run/docker.sock -d -e  AZUREUSERNAME -e AZUREPASSWORD -e SUBID -e LOCATION -e TEAMNAME -e RECIPIENTEMAIL -e CHATCONNECTIONSTRING -e CHATMESSAGEQUEUE -e TENANTID -e APPID -e GITBRANCH devopsoh/proctor-container &'
#echo "############### End of custom script ###############"