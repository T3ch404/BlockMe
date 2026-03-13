BlockMe is a suuuuuper simple project. However, your deployment can be as 
simple or complex as you'd like. We provide documentation and instructions to 
walk you through hosting BlockMe locally, in [Docker](docker.md), or a 
cloud-hosted [Kubernetes](kubernetes.md) environment.

## Deployment Considerations

### X-Forwarded-For Header
BlockMe was designed to operate behind a reverse proxy server for added 
security. With that in mind, BlockMe by default will block the IP address that 
is included in the X-Forwarded-For (XFF) header. If the XFF header is missing, 
BlockMe will default back to blocking the requester's IP. In the event that you 
forget to set up the XFF header in your reverse proxy's settings, BlockMe could
block your reverse proxy server resulting in all requests being blocked. To 
avoid this, you can add your reverse proxy's IP to the IGNORE_LIST environment 
variable. 