# WireGuard Controller Client

Client side service for [WireGuard Controller](https://github.com/wg-controller/wg-controller)

## Features

- Unattended enrollment with WireGuard Controller
- Automation of client IP routing, and NAT rules
- Synchronization of WireGuard keys and parameters with WireGuard controller
- Simplicity and performance benefits of [wireguard-go](https://github.com/WireGuard/wireguard-go)
- Single binary with minimal dependencies

## Installation

Download appropriate binary for your system

```
wget https://github.com/wg-controller/wg-controller-client/releases/download/latest/wg-controller-linux -O wg-controller
```

Make binary executable

```
sudo chmod +x wg-controller
```

Option 1: Run Standalone

```
sudo ./wg-controller --server-host wg.example.com --api-key kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU=
```

Option 2: Install as a service

```
sudo ./wg-controller --server-host wg.example.com --api-key kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU= --install && \
systemctl enable wg-controller && \
systemctl start wg-controller
```

## Command Line Options

| Flag           | Default  | Example                                      | Description                      |
| -------------- | -------- | -------------------------------------------- | -------------------------------- |
| --wg-interface | wg0      | utun11                                       | name used for kernel interface   |
| --server-host  | required | wg.example.com                               | public endpoint of wg-controller |
| --server-port  | 443      | 3000                                         | public port of wg-controller     |
| --api-key      | required | kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU= | api key created on wg-controller |
| --install      | false    |                                              | installs system service files    |
| --uninstall    | false    |                                              | cleans up system service files   |
