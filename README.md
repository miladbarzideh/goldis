# Build Your Own Redis with Go

[![License](https://img.shields.io/badge/license-apache_2.0-red?style=for-the-badge&logo=none)](LICENSE)

This project is a practical implementation of the concepts presented in the book ["Build Your Own Redis with C/C++"](https://build-your-own.org/redis/). The goal of this project is to provide a hands-on learning experience by building a simplified version of the Redis-like server using Go.

## Table of Contents

- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Concepts Explored](#concepts-explored)
- [Contributing](#contributing)

## Introduction

This project is a result of studying the book "Build Your Own Redis with C/C++". It aims to help developers gain a deeper understanding 
of network programming and data structures.

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository: `git clone https://github.com/miladbarzideh/goldis.git`
2. Navigate to the project directory: `cd goldis`
3. Run the application: `cd cmd/goldis/ && go build && ./goldis`
4. Use netcat to communicate with server: `nc localhost 6380`
5. Apply any basic command like: `set key value`

## Server Commands

1. SET: `SET key value`
2. GET: `GET key`
3. DEL: `DEL key`
4. KEYS: `KEYS`
5. PEXPIRE: `PEXPIRE key 10000` (ms)
6. PTTL: `PTTL key`
7. ZADD: `ZADD key 20 name`
8. ZSCORE: `ZSCORE key name`
9. ZREM: `ZREM key name`
10. ZQUERY: `ZQUERY key 18 name 0 10`
11. ZSHOW: `ZSHOW key`

## Concepts Explored

Throughout the development of this project, the following key concepts were explored and implemented:

| Concepts Explored                 |                       Implemented Features                        |    Further Steps |
|-----------------------------------|:-----------------------------------------------------------------:|-----------------:|
| Network programming               |                    Nonblocking IO, Event loop                     | Protocol parsing |
| Hashtable                         |            Hashtable, Chaining, Resizing, Intrusive DS            |                  |
| AVL Tree                          |                           Intrusive DS                            |                  |
| Sorted Set                        |                       Hashtable + AVL Tree                        |        Skip List |
| Timers                            |                     Kick out idle connections                     |                  |
| Heap and TTL                      |                         TTL with Min Heap                         |                  |
| Thread Pool - Asynchronous Tasks  | The producer-consumer problem, Synchronization primitives (Mutex) |   Try other ways |

## Contributing

Contributions to this project are welcome! If you find any issues or want to add new features, feel free to create a pull request.
