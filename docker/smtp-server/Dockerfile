FROM python:3.11.1-alpine3.17

LABEL Description="Dummy SMTP server"

# Parameters for default user:group
ARG uid=1000
ARG user=appuser
ARG gid=1000
ARG group=appgroup

# Add user and group for build and runtime
RUN addgroup -g "${gid}" "${group}" && adduser -D -h /home/${user} -s /bin/bash -G "${group}" -u "${uid}" "${user}"

# Prepare directories
RUN DIRS="/app" && \
    mkdir -p ${DIRS} && \
    chown -R ${user}:${group} $DIRS

# Switch to non-root user
USER ${user}

# Upgrade Pip from PyPi
ARG PIP_VERSION='22.3.1'
RUN python3 -m pip install "pip==${PIP_VERSION}"

# Install requirements
COPY ./requirements.txt /app/requirements.txt
RUN python3 -m pip install --user --upgrade -r /app/requirements.txt

WORKDIR /app
EXPOSE 1025

# Start SMTP server by default
CMD python3 -m aiosmtpd --nosetuid --listen 0.0.0.0:1025
