## Manual Installation
* build: `GOOS=linux go build -o gome-schedule .`
* service file goes into `/usr/lib/systemd/system/`
* env file goes into `/etc/gome/schedule/`
* database is located `/etc/gome/schedule/db`
* reload systemd: `systemctl daemon-reload`
* open FW: `firewall-cmd --zone=public --permanent --add-port 6661/tcp`
* save FW: `firewall-cmd --reload`