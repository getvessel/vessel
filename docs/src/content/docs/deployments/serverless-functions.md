---
title: Serverless Functions
description: Write and deploy serverless functions directly from the dashboard.
---

Vessl supports serverless function deployment for lightweight, event-driven workloads.

## Supported Runtimes

- **Node.js** — ECMAScript modules, CommonJS
- **Python** — Python 3.x
- **Go** — Go modules

## Creating a Serverless Function

1. Go to your project's **Services** tab.
2. Click **New Service → Serverless Function**.
3. Select a runtime (Node.js, Python, or Go).
4. Write your function code in the in-browser editor.
5. Click **Save & Deploy**.

The function is immediately deployed and available at a URL within your project's domain.

## In-Browser Editor

Vessl includes a Monaco-based code editor directly in the dashboard:

- Syntax highlighting for all supported runtimes
- Multi-file support
- Save and deploy with a single click
- Version history

## Function Structure

### Node.js

```js
export function handle(request) {
  return new Response('Hello from Vessl!', {
    headers: { 'content-type': 'text/plain' },
  });
}
```

### Python

```python
def handle(request):
    return {
        'statusCode': 200,
        'body': 'Hello from Vessl!',
        'headers': {
            'content-type': 'text/plain'
        }
    }
```

### Go

```go
package main

import (
    "io"
    "net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello from Vessl!")
}
```

## Environment Variables

Serverless functions have access to the same [environment variables](/deployment/#environment-variables) as your services, including auto-linked database connection strings.

## Deployment

Functions are deployed as lightweight containers with fast cold-start times. Deployments are triggered:

- When you click **Save & Deploy** from the editor
- Via the API for CI/CD integration

## Logs

View function logs from the **Logs** tab in the service detail page. Logs are streamed in real-time and persisted for 7 days.
