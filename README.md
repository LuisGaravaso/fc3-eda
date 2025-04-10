
# 💸 Microservices Wallet Project

A practical implementation of a **microservices architecture** with **event-driven communication** using **Kafka**, developed in **Go**. 

This project simulates a digital wallet system composed of two core services:

- **Wallet Service**: Handles clients, accounts, and transactions.
- **Balance Service**: Listens to transaction events and keeps account balances updated in real-time.

---

## 🧠 Architecture Highlights

- **Event-Driven Design**: Transaction events are published by the Wallet Service and consumed by the Balance Service.
- **Service Autonomy**: Each microservice has its own isolated database.
- **RESTful APIs**: Services expose simple HTTP interfaces for integration and testing.

---

## 🧩 Services Overview

### 🚀 Wallet Service

- **Port**: `8080`
- Manages clients, accounts, and transactions.
- Implements business logic for transfers.
- Publishes `BalanceUpdated` events to Kafka.

**Available Endpoints**:
| Method | Endpoint             | Description                      |
|--------|----------------------|----------------------------------|
| POST   | `/clients`           | Create a new client              |
| POST   | `/accounts`          | Create a new account             |
| POST   | `/transactions`      | Perform a transaction            |
| GET    | `/health`            | Health check                     |

---

### 📊 Balance Service

- **Port**: `3003`
- Maintains a **read-optimized view** of balances.
- Subscribes to Kafka to receive balance updates.

**Available Endpoints**:
| Method | Endpoint                      | Description                        |
|--------|-------------------------------|------------------------------------|
| GET    | `/balances/{account_id}`      | Retrieve current balance for account |
| GET    | `/health`                     | Health check                       |

---

## 🛠️ Getting Started

### Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Run the Full Stack

Clone the repository:

```bash
git clone https://github.com/LuisGaravaso/fc3-eda.git
cd fc3-eda
```

Start all services:

```bash
docker-compose up -d
```

This will start:

- MySQL (for both services)
- Kafka & Zookeeper
- Wallet Service
- Balance Service

All components will:

- Run DB migrations automatically
- Seed sample data

---

## 🔄 How to Interact with the API

You can use the included `client.http` file or use `cURL` to interact with the services.

### Option 1: VS Code (Recommended)

1. Install the **REST Client** extension by *Huachao Mao*
2. Open `api/client.http`
3. Click **"Send Request"** above each command

### Option 2: Use cURL manually

```bash
# Get balance - should be 100
curl http://localhost:3003/balances/7ebc23f5-dd1e-4d93-9490-9fce5052a5f5

# Get balance - should be 100
curl http://localhost:3003/balances/dff2d137-bba6-4138-81b9-3da7567f122b

# Transfer 10 from A to B
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"account_id_from": "7ebc23f5-dd1e-4d93-9490-9fce5052a5f5", "account_id_to": "dff2d137-bba6-4138-81b9-3da7567f122b", "amount": 10}'

# Get balances after transfer
curl http://localhost:3003/balances/7ebc23f5-dd1e-4d93-9490-9fce5052a5f5  # Should be 90
curl http://localhost:3003/balances/dff2d137-bba6-4138-81b9-3da7567f122b  # Should be 110
```

---

## 🔁 Data Flow Explained

1. A client makes a transaction request to the Wallet Service.
2. The transaction is processed and persisted.
3. A `BalanceUpdated` event is published to Kafka.
4. The Balance Service consumes the event and updates its local balance view.
5. Clients can query balances via the Balance Service at any time.

---

## ⚙️ Technical Details

- Wallet Service implements the **Unit of Work** pattern for transaction integrity.
- Balance Service uses **Kafka event handlers** to update balances.
- Health endpoints are provided for both services.
- Database schemas and sample data are initialized automatically at startup.

---

## 🧰 Tech Stack

| Component         | Technology     |
|------------------|----------------|
| Language         | Go             |
| Messaging        | Kafka          |
| Database         | MySQL          |
| Containerization | Docker Compose |
