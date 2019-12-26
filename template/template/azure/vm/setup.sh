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
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates p11-kit
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y curl wget
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y gnupg-agent
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y software-properties-common
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y docker-ce docker-ce-cli containerd.io
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y unzip

#Add user to docker usergroup
sudo DEBIAN_FRONTEND=noninteractive apt-get remove -y unscd
sudo usermod -aG docker azureuser

#Holding walinuxagent before upgrade
sudo DEBIAN_FRONTEND=noninteractive apt-mark hold walinuxagent
sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y

# Install JDK
# Inspired by https://raw.githubusercontent.com/docker-library/openjdk/master/14/jdk/slim/Dockerfile
export LANG=C.UTF-8
export JAVA_HOME=/usr/java/openjdk-14
export PATH=$JAVA_HOME/bin:$PATH
export JAVA_VERSION=14-ea+26
export JAVA_URL=https://download.java.net/java/early_access/jdk14/26/GPL/openjdk-14-ea+26_linux-x64_bin.tar.gz
export JAVA_SHA256=abe716bf202ab8afe33e422bf83d05743def6a08b3b9843339cde74d1690e7cc

set -eux; \
	\
	savedAptMark="$(apt-mark showmanual)"; \
	wget -O openjdk.tgz "$JAVA_URL"; \
	echo "$JAVA_SHA256 openjdk.tgz" | sha256sum -c -; \
	\
	mkdir -p "$JAVA_HOME"; \
	tar --extract \
		--file openjdk.tgz \
		--directory "$JAVA_HOME" \
		--strip-components 1 \
		--no-same-owner \
	; \
	rm openjdk.tgz; \
	\
	apt-mark auto '.*' > /dev/null; \
	[ -z "$savedAptMark" ] || apt-mark manual $savedAptMark > /dev/null; \
	apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
	\
# update "cacerts" bundle to use Debian's CA certificates (and make sure it stays up-to-date with changes to Debian's store)
# see https://github.com/docker-library/openjdk/issues/327
#     http://rabexc.org/posts/certificates-not-working-java#comment-4099504075
#     https://salsa.debian.org/java-team/ca-certificates-java/blob/3e51a84e9104823319abeb31f880580e46f45a98/debian/jks-keystore.hook.in
#     https://git.alpinelinux.org/aports/tree/community/java-cacerts/APKBUILD?id=761af65f38b4570093461e6546dcf6b179d2b624#n29
	{ \
		echo '#!/usr/bin/env bash'; \
		echo 'set -Eeuo pipefail'; \
		echo 'if ! [ -d "$JAVA_HOME" ]; then echo >&2 "error: missing JAVA_HOME environment variable"; exit 1; fi'; \
# 8-jdk uses "$JAVA_HOME/jre/lib/security/cacerts" and 8-jre and 11+ uses "$JAVA_HOME/lib/security/cacerts" directly (no "jre" directory)
		echo 'cacertsFile=; for f in "$JAVA_HOME/lib/security/cacerts" "$JAVA_HOME/jre/lib/security/cacerts"; do if [ -e "$f" ]; then cacertsFile="$f"; break; fi; done'; \
		echo 'if [ -z "$cacertsFile" ] || ! [ -f "$cacertsFile" ]; then echo >&2 "error: failed to find cacerts file in $JAVA_HOME"; exit 1; fi'; \
		echo 'trust extract --overwrite --format=java-cacerts --filter=ca-anchors --purpose=server-auth "$cacertsFile"'; \
	} > /etc/ca-certificates/update.d/docker-openjdk; \
	chmod +x /etc/ca-certificates/update.d/docker-openjdk; \
	/etc/ca-certificates/update.d/docker-openjdk; \
	\
# https://github.com/docker-library/openjdk/issues/331#issuecomment-498834472
#	find "$JAVA_HOME/lib" -name '*.so' -exec dirname '{}' ';' | sort -u > /etc/ld.so.conf.d/docker-openjdk.conf; \
#	ldconfig; \
#	\
# https://github.com/docker-library/openjdk/issues/212#issuecomment-420979840
# https://openjdk.java.net/jeps/341
#	java -Xshare:dump; \
#	\
# basic smoke test
	javac --version; \
	java --version

# Install JMeter
# Inspired by this image https://hub.docker.com/r/cirit/jmeter

export JMETER_VERSION=5.2
export JMETER_HOME=/home/azureuser/apache-jmeter-${JMETER_VERSION}
export PATH=${JMETER_HOME}/bin:${PATH}
export CUSTOM_PLUGIN_VERSION=2.9

wget http://www-us.apache.org/dist/jmeter/binaries/apache-jmeter-${JMETER_VERSION}.tgz && \
	tar -xzf apache-jmeter-${JMETER_VERSION}.tgz -C /home/azureuser/
wget https://jmeter-plugins.org/files/packages/jpgc-casutg-${CUSTOM_PLUGIN_VERSION}.zip
unzip -o jpgc-casutg-${CUSTOM_PLUGIN_VERSION}.zip -d ${JMETER_HOME}

# Disable Rmi SSL option
sed -i.bak -e "s/#server.rmi.ssl.disable=false/server.rmi.ssl.disable=true/" ${JMETER_HOME}/bin/jmeter.properties

# Clean up

sudo DEBIAN_FRONTEND=noninteractive rm -rf apache-jmeter-${JMETER_VERSION}.tgz \
            jpgc-casutg-${CUSTOM_PLUGIN_VERSION}.zip \
			${JMETER_HOME}/bin/examples \
			${JMETER_HOME}/bin/templates \
			${JMETER_HOME}/bin/*.cmd \
			${JMETER_HOME}/bin/*.bat \
			${JMETER_HOME}/docs \
			${JMETER_HOME}/printable_docs && \
	apt-get -y remove wget && \
	apt-get -y --purge autoremove && \
	apt-get -y clean && \
	rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# log directory creation TODO Delete if it is not used.
mkdir -p /home/azureuser/logs

# Adding environment variables setup to azureuser
echo "# JMeter setup" >> /home/azureuser/.bashrc
echo "export JAVA_HOME=${JAVA_HOME}" >> /home/azureuser/.bashrc
echo "export PATH=${JMETER_HOME}/bin:${PATH}" >> /home/azureuser/.bashrc

# change owner of JMeter
chown -R azureuser $JMETER_HOME
chgrp -R azureuser $JMETER_HOME

cd /home/azureuser
# install volley
GET_VOLLEY=get_volley.sh
curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/${GIT_BRANCH}/script/get_volley.sh -o $GET_VOLLEY
chmod +x $GET_VOLLEY
chown azureuser $GET_VOLLEY
chgrp azureuser $GET_VOLLEY
sudo -u azureuser ./get_volley.sh

VOLLEY_START_SCRIPT=/home/azureuser/start_volley.sh
curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/${GIT_BRANCH}/script/start_volley.sh -o $VOLLEY_START_SCRIPT
chmod +x $VOLLEY_START_SCRIPT
chown azureuser $VOLLEY_START_SCRIPT
chgrp azureuser $VOLLEY_START_SCRIPT
# Master: 
# Start volley server
# Add cron for enabling start volley server when it starts

sudo -u azureuser --preserve-env=PATH,JAVA_HOME $VOLLEY_START_SCRIPT
echo "PATH=${PATH}" | crontab -u azureuser -
(crontab -l 2>/dev/null; echo "@reboot ${VOLLEY_START_SCRIPT}") | crontab -u azureuser -

# Slave:
# Start JMeter server
# Add cron for enabing start JMeter as server when it starts
# JMETER_SLAVE_START_SCRIPT=/home/azureuser/start_jmeter_slave.sh
# curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/master/script/start_jmeter_slave.sh -o $JMETER_SLAVE_START_SCRIPT
# chmod +x $JMETER_SLAVE_START_SCRIPT
# chown azureuser $JMETER_SLAVE_START_SCRIPT
# chgrp azureuser $JMETER_SLAVE_START_SCRIPT
# echo "@reboot ${JMETER_SLAVE_START_SCRIPT}" | crontab -u azureuser -


# /bin/bash -c 'docker run -d --name docker-daemon --privileged docker:stable-dind &'
# /bin/bash -c 'docker run -v /home/nginx/config:/home/nginx/config -v /home/nginx/contents:/home/nginx/contents -v /home/azureuser/logs:/home/azureuser/logs -v /var/run/docker.sock:/var/run/docker.sock -d -e  AZUREUSERNAME -e AZUREPASSWORD -e SUBID -e LOCATION -e TEAMNAME -e RECIPIENTEMAIL -e CHATCONNECTIONSTRING -e CHATMESSAGEQUEUE -e TENANTID -e APPID -e GITBRANCH devopsoh/proctor-container &'
#echo "############### End of custom script ###############"
