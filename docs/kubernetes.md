# Run BlockMe In Kubernetes (Advanced)
It's impossible to account for every Kubernetes environment, but we've done our 
best to put together a Kubernetes definition files that make hosting BlockMe in 
Kubernetes as easy as possible. 

## Configurations

### Service
In its default configuration, `kubernetes.yaml` exposes the BlockMe
service through a ClusterIP service. This works well in environments with
existing ingress controllers or CloudFlare tunnel pods.

There is also a NodePort configuration that has been commented out. This is a 
good option if you are hosting on a bear metal cluster and want to port-forward.

### Environment Variables
An example ConfigMap is available in configMap.yaml file. This ConfigMap can be 
modified using the configurations defined in [configuration](configuration.md)

## Deploy
1. Clone the BlockMe repository
    - `git clone https://github.com/T3ch404/BlockMe.git && cd BlockMe`
2. Modify kubernetes.yaml & configmap.yaml to meet your needs
3. Create the blockme namespace
    - `kubectl create namespace blockme`
3. Apply the Kubernetes definition files
    - `kubectl apply -f kubernetes`
