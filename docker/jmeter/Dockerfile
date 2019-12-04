FROM openjdk:14-jdk-slim 

# Inspired by this image https://hub.docker.com/r/cirit/jmeter

ENV JMETER_VERSION=5.2
ENV JMETER_HOME=/usr/local/apache-jmeter-${JMETER_VERSION}
ENV PATH=${JMETER_HOME}/bin:${PATH}

# Heap size settings if necessary. 
# For more details https://jmeter.apache.org/usermanual/get-started.html
# Usually these parameters need to configured
# ENV JVM_ARGS="-Xms1024m -Xmx1024m"
# ENV "-Xms1G -Xmx1G -XMaxMetaspaceSize=192m

RUN apt-get -y update && apt-get -y install wget

RUN wget http://www-us.apache.org/dist/jmeter/binaries/apache-jmeter-${JMETER_VERSION}.tgz && \
	tar -xzf apache-jmeter-${JMETER_VERSION}.tgz -C /usr/local/

RUN rm -rf apache-jmeter-${JMETER_VERSION}.tgz \
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