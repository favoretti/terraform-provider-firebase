---
page_title: "firestore_document Resource"
description: |-
  Manages a Firestore document.
---

# firestore_document (Resource)

Manages a Firestore document with full CRUD operations.

## Example Usage

### Basic Document

```hcl
resource "firestore_document" "user" {
  collection  = "users"
  document_id = "user-123"
  fields = jsonencode({
    name  = "John Doe"
    email = "john@example.com"
    age   = 30
  })
}
```

### Document with Complex Types

```hcl
resource "firestore_document" "user" {
  collection  = "users"
  document_id = "user-456"
  fields = jsonencode({
    name     = "Jane Doe"
    active   = true
    score    = 95.5
    tags     = ["admin", "developer"]
    metadata = {
      created_by = "terraform"
      version    = 1
    }
  })
}
```

### Subcollection Document

```hcl
resource "firestore_document" "order" {
  collection  = "users/user-123/orders"
  document_id = "order-001"
  fields = jsonencode({
    product  = "Widget"
    quantity = 5
  })
}
```

### Auto-generated Document ID

```hcl
resource "firestore_document" "log_entry" {
  collection = "logs"
  fields = jsonencode({
    message = "Application started"
    level   = "info"
  })
}
```

## Schema

### Required

- `collection` (String) - The collection path (e.g., "users" or "users/123/orders").
- `fields` (String) - JSON string of document fields.

### Optional

- `document_id` (String) - The document ID. If not provided, one will be auto-generated.
- `project` (String) - The GCP project ID. Overrides the provider project.
- `database` (String) - The Firestore database ID. Overrides the provider database.

### Read-Only

- `name` (String) - The full document resource name.
- `create_time` (String) - The time the document was created.
- `update_time` (String) - The time the document was last updated.

## Import

Documents can be imported using the full path or short path format:

```bash
# Full format: project/database/collection/document_id
terraform import firestore_document.example my-project/(default)/users/user-123

# Short format (uses provider defaults): collection/document_id
terraform import firestore_document.example users/user-123

# Subcollection format
terraform import firestore_document.example users/user-123/orders/order-001
```
