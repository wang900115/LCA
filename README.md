
 <h1 align="center"><strong>Lightweight Chat Application</strong></h1> 
    <p align="center">
      <img src="assets/banner.png" alt="LCA Banner" width="300"/>
    </p>
  
## Overview
  **LCA** is a message chat hybrid system written in Go. It is designed for Decentralized encrypted communication. However, central management is still the priority for system security.

## Knowledge
  - **Decentral Architecture:** The core idea of this project is decentralization, inspired by peer-to-peer gossip protocols. This approach naturally leads to the formation of strong, resilient communities. Each node acts not only as a giver but also as a receiver, enabling a fully participatory network. The interaction is staggered between the Opher-Chain and external blockchain ecosystems (such as Ethereum and others). This design promotes interoperability while maintaining independence.
  
  - **Encryption:** Messages are encrypted on the sender side using the RSA algorithm. A Merkle Patricia Tree is employed to ensure complete integrity and verification of nodes within the network.
  
  - **Distributed System:** The system is designed following distributed system principles, including the CAP theorem (Consistency, Availability, Partition Tolerance) and BASE (Basically Available, Soft State, Eventual Consistency). Quorum-based system design is applied to implement consensus algorithms such as Paxos and Raft, ensuring reliability and consistency across nodes.
  
  - **Security:** The project uses PASETO and JWT technologies for authentication and session management, enhancing system security and making unauthorized access or hacking significantly more difficult. 

## Prerequisite 
  - Golang Version >= 1.25.0
  - Docker Installed (images come from pulling hub)
  - If using Local (should have postgresql server and redis service)

## Download Build 

You can download the lastest build here:
  - [Windows](https://github.com/wang900115/LCA/releases/latest/download/main.exe)
  - [Linux](https://github.com/wang900115/LCA/releases/latest/download/main)

## Get Started
> [!WARNING] if Using local please check prerequsite 
  - *Docker*
    -  Run:  `docker-compose up --build`
    -  ShutDowan:  `docker-compose down`
  -  *Local* 
    - Window: 
      -  Build: `go build -o build ./cmd/LCA/main.go`
      -  Run: `./build/main.exe`
    - Linux:
      - Build: `make build`
      - Run: `make run`

## Question
  If you have any questions, please send me the ISSUE. I will personally understand and check if there are any omissions. Keep doing the best.

## Contributer
  - Main Dev: 
    - Name: Perry
## Licensing
  This project, LightWeight Chat Application (LCA), is released under an open-source license to encourage collaboration, transparency, and innovation in decentralized secure communication systems. We currently use the following license: MIT License You are free to: Use, Copy, Modify, Merge, Publish, and Distribute the software Use it for personal, educational, or commercial purposes Provided that: You include the original copyright and license You provide attribution to the original authors For the full license text, refer to the LICENSE file in the repository.