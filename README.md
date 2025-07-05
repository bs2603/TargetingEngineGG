# TargetingEngineGG

This repository implements a targeting engine that routes campaigns to the right requests, with a focus on **extensible, maintainable, and reusable code**.

Let's LLD right into it. 🙂

---

## Clean Coding and LLD

Clean coding principles like **SOLID** were popularized by *Uncle Bob*. Traditionally, Low-Level Design (LLD) examples have been illustrated in Java, an OOP-focused language with extensive support for inheritance and related constructs.

Unlike Java, **Go** is a statically typed language designed with **concurrency** in mind:
- It does not have classes or classical inheritance.
- It provides OOP-like features via **structs, methods, and interfaces**.
- Go favors **composition over inheritance**, organizing code around packages and functions.

---

## Database Setup

1. **Start MySQL**

   ```bash
   brew services start mysql

2. **Create the database and tables**

   ```bash
   mysql -u root -p < schema.sql

---

## Scalable Targeting Logic

Targeting logic and dimension handling are decoupled, so you can add new dimensions easily without changing application code.

- Dimension is stored as a VARCHAR, not an ENUM.

    - ENUM requires a schema change every time you add a new dimension.

    - VARCHAR allows arbitrary dimension names without DB migration.

    ```bash
    INSERT INTO targeting_rules (campaign_id, dimension, type, value)
    VALUES ('spotify', 'DeviceType', 'INCLUDE', 'Tablet');

---

## Targeting Dimensions

| Dimension          | Required | Example Value     |
| ------------------ | -------- | ----------------- |
| `app`              | ✅        | `com.example.app` |
| `country`          | ✅        | `US`              |
| `os`               | ✅        | `android`         |
| `devicetype`       | Optional | `tablet`          |
| `subscriptiontier` | Optional | `premium`         |


--- 

## Adding a New dimension

1. Update the RequiredDimensions or OptionalDimensions list in delivery/config.go.
2. Add the corresponding targeting rules in the DB:

    ```bash
    INSERT INTO targeting_rules (campaign_id, dimension, type, value)
    VALUES ('spotify', 'SubscriptionTier', 'INCLUDE', 'Premium');

--- 

## Handling the read heavy workload

1. In a service like this, if there are ~1000s of campaigns, thats KBs of data. And with 100000 of requests, thats still MBs of data. Adding an In-memory cache using the service's memory would be the fastest way to cache, as it will avoid the network hop required to query redis or anything.

2. But Since I am trying to make a production ready service, such a service can be deployed on multiple pods, which can be sharing cache. Its better to implement redis, which is what I will do :)

###Latency:
Redis fetches are ~0.2ms (microseconds), much faster than MySQL but a bit slower than in-process cache.

###Consistency:
Redis is always in sync every 30s.

###Scalability:
Multiple Go processes can all use the same Redis.

###Persistence:
Redis can persist to disk if you enable RDB/AOF.

On startup:
Loads all active campaigns and their rules from MySQL.
Store them in Redis as JSON blobs.

Every N seconds (e.g., 30s):
Refresh all campaigns in Redis, so the cache stays up to date with MySQL changes (new campaigns, new targeting rules, status changes).

On each request:

Get all campaigns from Redis (or fetch per campaign if you prefer).

Match campaigns in memory.

Return results.

Cache fallback:

If Redis is unavailable, optionally fall back to MySQL.

2. add code for this -
Fallback to MySQL
If Redis returns an error or no keys, you can load from MySQL as backup.
 This is useful if Redis is down.