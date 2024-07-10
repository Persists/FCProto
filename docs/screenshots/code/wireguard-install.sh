sudo apt update
sudo apt install wireguard

wg genkey | tee privatekey | wg pubkey > publickey

sudo wg-quick down wg0
sudo wg-quick up wg0
