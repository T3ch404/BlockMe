# Docker

## Quick Start
The fastest way to run BlockMe is by pulling the public pre-built image and 
running it in Docker locally.

`docker run -p 8080:8080 ghcr.io/t3ch404/blockme/blockme:latest`

## Run BlockMe Using Docker Compose
If you prefer to keep your services defined as code (you should), then Docker 
Compose is the next step in the right direction towards that goal. The only 
additional requirement beyond the Docker runtime environment is the installation 
of the Docker Compose plugin 
[Docker Compose Installation](https://docs.docker.com/compose/install/)

1. Clone the BlockMe repository
    - `git clone https://github.com/T3ch404/BlockMe.git && cd BlockMe`
2. Run the Docker Compose file
   - `docker-compose up` (Add `-d` to the end of this command to run in the 
background)

## Build the Docker Image
You can build the docker image yourself using the following commands:

- `git clone https://github.com/T3ch404/BlockMe.git && cd BlockMe`
- `docker build -t block-me:local .`

\
\
\
Hosing with Kubernetes? Check out the [Kubernetes Documentation](kubernetes.md)