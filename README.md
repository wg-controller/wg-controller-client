# WireGuard Controller Client

Client side service for [WireGuard Controller](https://github.com/wg-controller/wg-controller)

## Features

- Unattended enrollment with WireGuard Controller
- Automation of ip routing, name servers and nat rules on client
- Simplicity and performance benefits of [wireguard-go](https://github.com/WireGuard/wireguard-go)
- Minimal dependencies (all in one binary)

## Installation

// To Do

## Command Line Options

| Flag           | Default  | Example                                      | Description                      |
| -------------- | -------- | -------------------------------------------- | -------------------------------- |
| --wg-interface | wg0      | utun11                                       | name used for kernel interface   |
| --server-host  | required | wg.example.com                               | public endpoint of wg-controller |
| --server-port  | 443      | 3000                                         | public port of wg-controller     |
| --api-key      | required | kZdMQsztB-vR6Wve2dYYUOf6LXl5n2cgeESN8i7MQkU= | api key created on wg-controller |
| --install      | false    |                                              | installs system service files    |
| --uninstall    | false    |                                              | cleans up system service files   |
