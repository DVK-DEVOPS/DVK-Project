
# Variables
variable "resource_group_name" {
  description = "Resource group name"
  default     = "rg-monitoring-prod"
}

variable "app_resource_group_name" {
  description = "App Resource group name"
  default     = "rg-monitoring-prod"
}

variable "location" {
  description = "Azure region"
  default     = "northeurope"
}

variable "admin_username" {
  description = "Admin username for VM"
  default     = "azureuser"
}

variable "ssh_public_key_path" {
  description = "Path to SSH public key"
  default     = "~/.ssh/id_rsa.pub"
}


variable "app_vm_private_ip" {
  description = "Private IP of application VM"
  type        = string
}

variable "app_vnet_id" {
  description = "VNet ID of application VM (for peering)"
  type        = string
}

variable "app_vnet_name" {
  description = "VNet name of application VM"
  type        = string
}

variable "grafana_vm_private_ip" {
  description = "Private IP of existing Grafana VM"
  type        = string
}