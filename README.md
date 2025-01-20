# WireGuard Controller Client

Client side service for [WireGuard Controller](https://github.com/wg-controller/wg-controller)

## Features

- Unattended enrollment with WireGuard Controller
- Automation of IP routing, and NAT rules on client
- Synchronization of WireGuard keys and parameters with controller
- Simplicity and performance benefits of [wireguard-go](https://github.com/WireGuard/wireguard-go)
- Minimal dependencies (all in one binary)

## Installation

Download appropriate binary for your system
Linux AMD64

```
wget https://github.com/wg-controller/wg-controller-client/releases/download/v0.0.2/wg-controller-linux
```

Linux ARM64

```
wget https://github.com/wg-controller/wg-controller-client/releases/download/v0.0.2/wg-controller-linuxarm64
```

Make binary executable

```
sudo chmod +x wg-controller-linux
```

Run Standalone

```
sudo ./wg-controller-linux --server-host wg.example.com --api-key kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU=
```

OR

Install as a service

```
sudo ./wg-controller-linux --server-host wg.example.com --api-key kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU= --install \
systemctl enable wg-controller \
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
