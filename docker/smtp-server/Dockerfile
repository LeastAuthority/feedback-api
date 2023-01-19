FROM python:3.11.1-alpine3.17

LABEL Description="Dummy SMTP server"

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
