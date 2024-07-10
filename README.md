# FCProto

This repository contains a prototype developed for the Fog Computing lecture. The project demonstrates the key principles and practical applications of fog computing, including distributed data processing and edge device coordination.

## Cmd 
The directory contains the main entry points for the cloud and edge applications.

## Deployment
Configuration files for Docker and Terraform to deploy the application.

### Docker
Start the cloud environment with the following command:
```bash
docker compose up cloud
```

Start the edge application with the following command:
```bash
docker compose up edge
```

### Terraform
To provision the cloud infrastructure with Terraform, navigate to the `terraform` directory and run the following commands:
```bash
terraform init
terraform validate
terraform plan
terraform apply
```
### Ansible & Kubernetes
We originally tried to deploy the application with microk8s and ansible to make the usage scenario as realistic as possible. Unfortunately, we did not manage to get microk8s to use the cloud node completely. 
It was only classified as ready, but no pods were released. It only worked locally to run the fog computing deployment with node affinity. We thought it was a pity to throw away the big effort and therefore decided to keep the code in the repository. 

## Docs
The docs directory contain documentation for the project like the demonstration video and the project description report.

<video width="568" height="320" controls>
  <source src="./docs/show-case-recording.mp4" type="video/mp4">
</video>

## Internal
### Cloud
The cloud directory contains logic that is only executed on the cloud side. This includes env configurations, database logic with entities and logic to start up server.

### Edge
The edge directory contains logic that is only executed on the edge side. This includes env configurations and logic to start up the edge.

### Shared
The shared directory contains logic that is shared between the cloud and edge side. This includes the message model, queue logic, connection logic and utility functions.

## Pkg
The pkg directory contains all sensor logic developed for the project.

### Sensor
Sensor interval can be adjusted in the sensor client. 
The sensor logic is responsible for generating data and run on the edge application. There are three different sensors implemented:
- **VirtualSensor**: Generates random temperature & humidity data.
- **MemorySensor**: Generates memory usage data.
- **CpuSensor**: Generates CPU usage data.