<h1 align="center"><strong>Lightweight Chat Application (LCA)</strong></h1>

<p align="center">
  <img src="assets/inside.png" alt="LCA Banner" height="225" width="370" />
</p>

---

## 📌 Overview
**LCA** (Lightweight Chat Application) is a secure and lightweight and hybrid decentralized communication system.
It supports RESTful APIs, WebSocket messaging, and RPC protocols, combining centralized management for security
and decentralized personal data for resilience and privacy.

## 🧠 Core Concepts & Features
- **Core Features**
  - **DID** — Decentralized Identifiers for user identity management
  - **DCC**  — Decentralized Communication Channel for peer-to-peer messaging
  - **External Interface** — Supports fetching and interacting with on-chain data

- **Security Architecture**
  - **Hybrid Encryption** using Curve25519 + AES
  - **PASETO** and **JWT** for external api authentication and session management  
  - **Tamper Resistance** — Protects against unauthorized access and data modification
  - **Integrity Checking**  — Validates message integrity via CRC/HMAC
---

## 🔧 Prerequisites
Before you start, make sure you have:

- **Golang** `>= 1.25.0`
- **Docker** (images will be pulled automatically from Docker Hub)
- **Local setup (optional)**  
  - PostgreSQL server  
  - Redis service  

---
## Download Build 

You can download the lastest build here:
  - [Windows](https://github.com/wang900115/LCA/releases/latest/download/main.exe)
  - [Linux](https://github.com/wang900115/LCA/releases/latest/download/main)

## Get Started
> [!WARNING]  
> If running locally, please verify you meet the prerequisites above. 
  - *Docker*
    -  Run:  `docker-compose up --build`
    -  ShutDowan:  `docker-compose down`
  - *Local* 
    - Window: 
      -  Build: `go build -o build ./cmd/LCA/main.go`
      -  Run: `./build/main.exe`
    - Linux:
      -  Build: `make build`
      -  Run: `make run`
## Brief Sample 
``` mermaid
graph TD
    A[Node A] -->|Sign X25519 PubKey with Ed25519| B[Node B]
    A -->|Sign X25519 PubKey with Ed25519| C[Node C]
    B -->|Verify Signature & Create Private Channel| A
    C -->|Verify Signature & Create Private Channel| A
    C -->|Relay Communication| B
```

## ❓ Question
  If you have any questions, please send me the ISSUE. I will personally understand and check if there are any omissions. Keep doing the best.

## 👨‍💻 Contributer
  - Main Dev: 
    - Name: Perry
    - Name: Aliz
## 📄 Licensing
  This project, LightWeight Chat Application (LCA), is released under an open-source license to encourage collaboration, transparency, and innovation in decentralized secure communication systems. We currently use the following license: MIT License You are free to: Use, Copy, Modify, Merge, Publish, and Distribute the software Use it for personal, educational, or commercial purposes Provided that: You include the original copyright and license You provide attribution to the original authors For the full license text, refer to the LICENSE file in the repository.