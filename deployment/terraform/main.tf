resource "google_compute_instance" "default" {
  name         = var.instance_name
  machine_type = "e2-small"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
      size  = 50
    }
  }

  network_interface {
    network = "default"

    access_config {
      # Ephemeral public IP
    }
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/ansible.pub")}"
  }

  tags = ["allow-microk8s-ports"]
}

resource "google_compute_firewall" "allow_microk8s_ingress" {
  name    = "allow-microk8s-ingress"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["16443", "25000", "12379", "10250", "19001"]
  }

  allow {
    protocol = "udp"
    ports    = ["4789"]
  }

  source_ranges = ["0.0.0.0/0"]  # Open to all IPs; modify this to restrict access as needed

  target_tags = ["allow-microk8s-ports"]
}

resource "google_compute_firewall" "allow_microk8s_egress" {
  name    = "allow-microk8s-egress"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  destination_ranges = ["0.0.0.0/0"]

  direction = "EGRESS"
  target_tags = ["allow-microk8s-ports"]
}
