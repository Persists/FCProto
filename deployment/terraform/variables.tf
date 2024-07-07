variable "project_id" {
  description = "The project ID to deploy to"
  type        = string
  default     = "fogcomputing-428415"
}

variable "region" {
  description = "The region to deploy resources"
  type        = string
  default     = "europe-west3"
}

variable "zone" {
  description = "The zone to deploy resources"
  type        = string
  default     = "europe-west3-a"
}

variable "instance_name" {
  description = "Name of the VM instance"
  type        = string
  default     = "cloud-node"
}
