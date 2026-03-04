---
page_title: "firestore_documents Data Source"
description: |-
  Lists Firestore documents in a collection.
---

# firestore_documents (Data Source)

Lists documents in a Firestore collection with optional filtering, ordering, and pagination.

## Example Usage

### List All Documents

```hcl
data "firestore_documents" "all_users" {
  collection = "users"
}

output "user_count" {
  value = length(data.firestore_documents.all_users.documents)
}
```

### Filter Documents

```hcl
data "firestore_documents" "active_users" {
  collection = "users"

  where {
    field    = "active"
    operator = "EQUAL"
    value    = jsonencode(true)
  }
}
```

### Multiple Filters

```hcl
data "firestore_documents" "premium_active_users" {
  collection = "users"

  where {
    field    = "active"
    operator = "EQUAL"
    value    = jsonencode(true)
  }

  where {
    field    = "tier"
    operator = "EQUAL"
    value    = jsonencode("premium")
  }
}
```

### Order and Limit

```hcl
data "firestore_documents" "recent_orders" {
  collection = "orders"

  order_by {
    field     = "created_at"
    direction = "DESCENDING"
  }

  limit = 10
}
```

### Complex Query

```hcl
data "firestore_documents" "filtered_users" {
  collection = "users"

  where {
    field    = "age"
    operator = "GREATER_THAN_OR_EQUAL"
    value    = jsonencode(18)
  }

  order_by {
    field     = "age"
    direction = "ASCENDING"
  }

  order_by {
    field     = "name"
    direction = "ASCENDING"
  }

  limit = 100
}
```

## Schema

### Required

- `collection` (String) - The collection path (e.g., "users" or "users/123/orders").

### Optional

- `project` (String) - The GCP project ID. Overrides the provider project.
- `database` (String) - The Firestore database ID. Overrides the provider database.
- `limit` (Number) - Maximum number of documents to return.

### Blocks

#### where (Optional)

Filter conditions for the query. Multiple `where` blocks are combined with AND.

- `field` (String, Required) - The field path to filter on.
- `operator` (String, Required) - The comparison operator. Valid values:
  - `EQUAL`
  - `NOT_EQUAL`
  - `LESS_THAN`
  - `LESS_THAN_OR_EQUAL`
  - `GREATER_THAN`
  - `GREATER_THAN_OR_EQUAL`
  - `ARRAY_CONTAINS`
  - `IN`
  - `ARRAY_CONTAINS_ANY`
  - `NOT_IN`
- `value` (String, Required) - The value to compare against, JSON encoded.

#### order_by (Optional)

Ordering for the query results.

- `field` (String, Required) - The field path to order by.
- `direction` (String, Optional) - The sort direction. Valid values: `ASCENDING` (default), `DESCENDING`.

### Read-Only

- `documents` (List of Object) - List of documents in the collection. Each document contains:
  - `document_id` (String) - The document ID.
  - `fields` (String) - JSON string of document fields.
  - `create_time` (String) - The time the document was created.
  - `update_time` (String) - The time the document was last updated.
