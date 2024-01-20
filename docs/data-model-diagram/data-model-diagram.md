```mermaid
---
title: Transaction - Data Model Diagram
---
erDiagram
    CUSTOMER ||--|{ ACCOUNT : has
    CUSTOMER {
        string id PK
        string name
        string username
        string createdAt
    }
    ACCOUNT ||--|{ TRANSACTION : has
    ACCOUNT {
        string accountId PK
        string userId 
        string createdAt
    }
    TRANSACTION {
        string id PK
        string origin "e.g. “desktop-web”, “mobile-android”,
“mobile-ios”"
        string userId 
        string accountId
        float amount
        string operationType "credit or debit"
        string createdAt
    }
```

Initially, this is a representation of a NoSQL data model diagram. This is why there is no explicitly mentioned foreign key in the tables

NoSQL was the choice as a high volume of transactions is expected, as it can scale horizontally very well. Since the transaction ID serves as the partition key with high cardinality, it helps avoid hot partitions.