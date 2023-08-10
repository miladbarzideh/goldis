# Build Your Own Redis with Go

[![License](https://img.shields.io/badge/license-apache_2.0-red?style=for-the-badge&logo=none)](LICENSE)

This project is a practical implementation of the concepts presented in the book ["Build Your Own Redis with C/C++"](https://build-your-own.org/redis/). The goal of this project is to provide a hands-on learning experience by building a simplified version of the Redis-like server using Go.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
- [Concepts Explored](#concepts-explored)
- [Contributing](#contributing)

## Introduction

This project is a result of studying the book "Build Your Own Redis with C/C++". It aims to help developers gain a deeper understanding 
of network programming and data structures.

## Features

- Networking support for client-server communication
- Basic key-value storage
- In-memory data storage
- Basic command processing (e.g., GET, SET, DEL)


## Getting Started

To get started with the project, follow these steps:

1. Clone the repository: `git clone https://github.com/miladbarzideh/goldis.git`
2. Navigate to the project directory: `cd goldis`
3. Run the application: `cd cmd/goldis/ && go build && ./goldis`
4. Use netcat to communicate with server: `nc localhost 6380`
5. Apply any basic command like: `set key value`

## Concepts Explored

Throughout the development of this project, the following key concepts were explored and implemented:

| Concepts Explored  |     Implemented Features      |    Further Steps |
|--------------------|:-----------------------------:|-----------------:|
| Network programming |  Nonblocking IO, Event loop   | Protocol parsing |
| Server commands    |         GET, SET, DEL         |                  |
| Hashtable          | Hashtable, Chaining, Resizing |                  |
| Data Serialization |              WIP              |                  |
| AVL Tree           |              WIP              |                  |
| Data Serialization |              WIP              |                  |
| Sorted Set         |              WIP              |                  |
| Heap and TTL       |              WIP              |                  |
| Thread Pool        |              WIP              |                  |

## Contributing

Contributions to this project are welcome! If you find any issues or want to add new features, feel free to create a pull request.
