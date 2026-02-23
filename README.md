# BlockMe
This is a dumb project that creates a web app with a single button. When pressed, the user's public IP is added to a block list. All subsequent requests from that address are responded to with a 403 response.

Try it out for yourself at:\
https://blockme.titanicswimteam.com

## Self-Host
NOTE: This project relies on the X-Forwarded-For header in order to work. Make sure you host this project behind a reverse proxy that supports the X-Forwarded-For header.

**Docker Run**\
`docker run -p 8080:8080 ghcr.io/t3ch404/blockme/blockme:latest`

**Docker Compose**\
`git clone https://github.com/T3ch404/BlockMe.git` \
`cd BlockMe/`\
`docker-compose up`

**Kubernetes**\
`git clone https://github.com/T3ch404/BlockMe.git` \
`cd BlockMe/`\
`kubectl apply -f kubernetes.yaml`

## Coming Soon (in no particular order)
- Better instructions for self-hosting
- Logging
- Persistent block list (Maybe wildly overengineered using a PostgreSQL db with replicas?)
- CSS styling
- Maybe a Twitter bot or something?