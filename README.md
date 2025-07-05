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


