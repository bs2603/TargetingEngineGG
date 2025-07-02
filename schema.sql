CREATE DATABASE targeting;

USE targeting;

CREATE TABLE campaigns (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255),
    image_url VARCHAR(255),
    cta VARCHAR(100),
    state ENUM('ACTIVE','INACTIVE') NOT NULL
);

CREATE TABLE targeting_rules(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    campaign_id VARCHAR(50),
    dimention ENUM('APP','COUNTRY','OS'),
    type ENUM('INCLUDE','EXCLUDE'),
    value VARCHAR(255),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);