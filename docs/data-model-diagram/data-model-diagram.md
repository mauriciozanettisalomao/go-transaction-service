```mermaid
---
title: Transaction - Data Model Diagram
---
erDiagram
    CUSTOMER ||--|{ TRANSACTION : has
    CUSTOMER {
        string id PK
        string name
        string username
        string createdAt
    }
    TRANSACTION {
        string userId PK
        string transactionId 
        string origin "e.g. “desktop-web”, “mobile-android”,
“mobile-ios”"
        float amount
        string operationType "credit or debit"
        string createdAt
    }
```

Initially, this is a representation of a NoSQL data model diagram. This is why there is no explicitly mentioned foreign key in the tables

NoSQL was the choice as a high volume of transactions per user is expected, as it can scale horizontally very well,  ensuring stable performance. Since the User ID and transaction ID serve as the partition key with high cardinality, it helps avoid hot partitions.