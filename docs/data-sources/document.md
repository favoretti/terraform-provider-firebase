---
page_title: "firestore_document Data Source"
description: |-
  Retrieves a single Firestore document.
---

# firestore_document (Data Source)

Retrieves a single Firestore document by collection and document ID.

## Example Usage

```hcl
data "firestore_document" "user" {
  collection  = "users"
  document_id = "user-123"
}

output "user_name" {
  value = jsondecode(data.firestore_document.user.fields).name
}
```

### Subcollection Document

```hcl
data "firestore_document" "order" {
  collection  = "users/user-123/orders"
  document_id = "order-001"
}
```

## Schema

### Required

- `collection` (String) - The collection path (e.g., "users" or "users/123/orders").
- `document_id` (String) - The document ID to retrieve.

### Optional

- `project` (String) - The GCP project ID. Overrides the provider project.
- `database` (String) - The Firestore database ID. Overrides the provider database.

### Read-Only

- `fields` (String) - JSON string of document fields.
- `name` (String) - The full document resource name.
- `create_time` (String) - The time the document was created.
- `update_time` (String) - The time the document was last updated.
