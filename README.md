# BlockMe
This is a dumb project that creates a web app with a single button. When 
pressed, the user's public IP is added to a block list. All subsequent requests 
from that address are responded to with a 403 response.

Try it out for yourself at:\
https://blockme.titanicswimteam.com

## Self-Host
Self-hosting instructions are available in the docs directory for deployments 
using [docker](docs/docker.md) or [kubernetes](docs/kubernetes.md). 
Additionally, [configuration](docs/configuration.md) and 
[other hosting considerations](docs/hosting.md) are also included in the docs 
directory as well.

## Coming Soon (in no particular order)
- ~~Better instructions for self-hosting~~
- Logging
- ~~Persistent block list (Maybe wildly overengineered using a PostgreSQL db with replicas?)~~
- CSS styling
- Maybe a Twitter bot or something?