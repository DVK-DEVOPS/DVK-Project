terraform {
  required_version = ">= 1.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# Resource group
data "azurerm_resource_group" "monitoring" {
  name     = var.resource_group_name
}

# Virtual Network
resource "azurerm_virtual_network" "monitoring" {
  name                = "vnet-monitoring"
  address_space       = ["10.1.0.0/16"]
  location            = data.azurerm_resource_group.monitoring.location
  resource_group_name = data.azurerm_resource_group.monitoring.name
}

resource "azurerm_subnet" "monitoring" {
  name                 = "subnet-monitoring"
  resource_group_name  = data.azurerm_resource_group.monitoring.name
  virtual_network_name = azurerm_virtual_network.monitoring.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "monitoring" {
  name                = "pip-monitoring"
  location            = data.azurerm_resource_group.monitoring.location
  resource_group_name = data.azurerm_resource_group.monitoring.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

# Network Security Group
resource "azurerm_network_security_group" "monitoring" {
  name                = "nsg-monitoring"
  location            = data.azurerm_resource_group.monitoring.location
  resource_group_name = data.azurerm_resource_group.monitoring.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Prometheus"
    priority                   = 1003
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "9090"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Loki"
    priority                   = 1004
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3100"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "Alloy"
    priority                   = 1005
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "12345"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "*"
  }
}

# NIC
resource "azurerm_network_interface" "monitoring" {
  name                = "nic-monitoring"
  location            = data.azurerm_resource_group.monitoring.location
  resource_group_name = data.azurerm_resource_group.monitoring.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.monitoring.id
    private_ip_address_allocation = "Static"
    public_ip_address_id          = azurerm_public_ip.monitoring.id
  }
}

resource "azurerm_network_interface_security_group_association" "monitoring" {
  network_interface_id      = azurerm_network_interface.monitoring.id
  network_security_group_id = azurerm_network_security_group.monitoring.id
}

# Cloud-init for installing Docker / Alloy / Prometheus / Loki
data "cloudinit_config" "monitoring" {
  gzip          = true
  base64_encode = true

  part {
    content_type = "text/cloud-config"
    content = templatefile("${path.module}/cloud-init.yaml", {
      app_vm_ip        = var.app_vm_private_ip
      monitoring_vm_ip = azurerm_network_interface.monitoring.private_ip_address
    })
  }
}

# VM – switched to cheapest possible size
resource "azurerm_linux_virtual_machine" "monitoring" {
  name                = "vm-monitoring"
  resource_group_name = data.azurerm_resource_group.monitoring.name
  location            = data.azurerm_resource_group.monitoring.location
  size                = "Standard_B1s"  # free-tier friendly

  admin_username      = var.admin_username
  custom_data         = data.cloudinit_config.monitoring.rendered

  network_interface_ids = [
    azurerm_network_interface.monitoring.id,
  ]

  admin_ssh_key {
    username   = var.admin_username
    public_key = file(var.ssh_public_key_path)
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 32
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }

  identity {
    type = "SystemAssigned"
  }
}

# VNet peering
resource "azurerm_virtual_network_peering" "monitoring_to_app" {
  name                      = "peer-monitoring-to-app"
  resource_group_name       = data.azurerm_resource_group.monitoring.name
  virtual_network_name      = azurerm_virtual_network.monitoring.name
  remote_virtual_network_id = var.app_vnet_id
  allow_forwarded_traffic   = true
}

resource "azurerm_virtual_network_peering" "app_to_monitoring" {
  name                      = "peer-app-to-monitoring"
  resource_group_name       = var.app_resource_group_name
  virtual_network_name      = var.app_vnet_name
  remote_virtual_network_id = azurerm_virtual_network.monitoring.id
  allow_forwarded_traffic   = true
}

output "monitoring_vm_public_ip" {
  value = azurerm_public_ip.monitoring.ip_address
}

output "monitoring_vm_private_ip" {
  value = azurerm_network_interface.monitoring.private_ip_address
}

output "ssh_command" {
  value = "ssh ${var.admin_username}@${azurerm_public_ip.monitoring.ip_address}"
}

output "prometheus_url" {
  value = "http://${azurerm_network_interface.monitoring.private_ip_address}:9090"
}

output "loki_url" {
  value = "http://${azurerm_network_interface.monitoring.private_ip_address}:3100"
}
