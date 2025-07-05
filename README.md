# TargetingEngineGG

This repository implements a targeting engine that routes campaigns to the right requests, with a focus on **extensible, maintainable, and reusable code**.
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
