provider "google" {
  credentials = file("service-account-key.json")
  project     = "fogcomputing-428415"
  region      = "europe-west3"  # Frankfurt
  zone        = "europe-west3-a"
}
