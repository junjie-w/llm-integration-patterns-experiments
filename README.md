# LLM Integration Patterns Experiments

## API Endpoints

### 1. Basic LLM Completion

**Endpoint**: `POST /api/support/basic-llm-completion`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "message": "How can I return a product that I purchased?"
}
```

**Example Response**
```json
{
  "reply": "To return a product, please check the retailer's return policy for instructions. Typically, you will need to provide proof of purchase and return the item in its original packaging within a specified time frame."
}
```
</details>

### 2. Knowledge RAG

**Endpoint**: `POST /api/support/knowledge-rag`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "message": "How do I reset my password?",
  "use_vector_search": true
}
```

**Example Response**
```json
{
    "reply": "To reset your password, click the 'Forgot Password' link on the login page. You will receive an email with a link to create a new password.",
    "sources": [
        {
            "id": "doc_3",
            "title": "Account Password Reset",
            "content": "To reset your password, click the 'Forgot Password' link on the login page. You will receive an email with a link to create a new password. The link expires after 24 hours. If you don't receive the email, please check your spam folder.",
            "tags": [
                "account",
                "password",
                "login"
            ]
        }
    ]
}
```

**Example Request**
```json
{
  "message": "What is your return policy?",
  "use_vector_search": true
}
```

**Example Response**
```json
{
    "reply": "Our return policy allows returns within 30 days of purchase with a receipt. Items must be in original condition with all packaging. Refunds are processed to the original payment method within 5-7 business days.",
    "sources": [
        {
            "id": "doc_1",
            "title": "Return Policy",
            "content": "Our return policy allows returns within 30 days of purchase with a receipt. Items must be in original condition with all packaging. Refunds are processed to the original payment method within 5-7 business days.",
            "tags": [
                "returns",
                "policy",
                "refunds"
            ]
        },
        {
            "id": "doc_4",
            "title": "Product Warranty",
            "content": "All electronics come with a 1-year manufacturer warranty covering defects in materials and workmanship. The warranty does not cover damage from misuse, accidents, or normal wear and tear. Extended warranties are available for purchase.",
            "tags": [
                "warranty",
                "guarantee",
                "repairs"
            ]
        },
        {
            "id": "doc_5",
            "title": "Membership Benefits",
            "content": "Premium members receive free shipping on all orders, early access to sales, exclusive discounts, and priority customer support. Membership costs $49.99 per year and can be canceled at any time.",
            "tags": [
                "membership",
                "premium",
                "benefits"
            ]
        }
    ]
}
```
</details>

### 3. Function Calling

**Endpoint**: `POST /api/support/function-calling`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "message": "What's the status of my order ORD-1234?"
}
```

**Example Response**
```json
{
  "reply": "I've checked your order ORD-1234. It's currently in the 'Shipped' status and is being delivered by FedEx. The estimated delivery date is May 26, 2025. You can track your package using the tracking number TRK772807.",
  "tool_calls": [
    {
      "name": "customer_info",
      "arguments": {
        "query_type": "shipping",
        "order_id": "ORD-1234"
      },
      "result": {
        "order_id": "ORD-1234",
        "status": "Shipped",
        "shipping_info": {
          "carrier": "FedEx",
          "estimated_delivery": "2025-05-26",
          "status": "In Transit",
          "tracking_number": "TRK772807"
        }
      }
    }
  ]
}
```
</details>

### 4. Reasoning Agent

**Endpoint**: `POST /api/support/reasoning-agent`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "message": "I ordered headphones last week (order ID: ORD-1234) but I just noticed they were shipped to my old address. I need to know if they were delivered yet, and if not, can I still change the delivery address?"
}
```

**Example Response**
```json
{
  "agent_id": "41149034-e7f1-456f-9cae-d4dc4b981e15",
  "answer": "The headphones with order ID ORD-1234 have been shipped via UPS with a tracking number TRK485533. The estimated delivery date is May 26, 2025, and the current status is \"Order Placed.\" The order status is \"Processing,\" which means it has not been delivered yet. Since the order is still in processing, there is a possibility to change the delivery address. I recommend contacting customer support as soon as possible to request a change in the delivery address for your order.",
  "complete": true,
  "steps": [
    {
      "type": "action",
      "content": "Using tool 'customer_info' with args: {\"order_id\":\"ORD-1234\",\"query_type\":\"shipping\"}"
    },
    {
      "type": "observation",
      "content": "{\n  \"customer_email\": \"customer@example.com\",\n  \"date\": \"2025-05-03\",\n  \"items\": 1,\n  \"order_id\": \"ORD-1234\",\n  \"shipping_info\": {\n    \"carrier\": \"UPS\",\n    \"estimated_delivery\": \"2025-05-26\",\n    \"status\": \"Order Placed\",\n    \"tracking_number\": \"TRK485533\"\n  },\n  \"status\": \"Shipped\",\n  \"total_amount\": \"$105.53\"\n}"
    },
    {
      "type": "thought",
      "content": "Thought: The headphones with order ID ORD-1234 have been shipped via UPS with a tracking number TRK485533. The estimated delivery date is May 26, 2025. The current status is \"Order Placed.\"\n\nAction: I will check if the delivery address can still be changed for this order."
    },
    {
      "type": "action",
      "content": "Using tool 'customer_info' with args: {\"order_id\":\"ORD-1234\",\"query_type\":\"order\"}"
    },
    {
      "type": "observation",
      "content": "{\n  \"customer_email\": \"customer@example.com\",\n  \"date\": \"2025-05-04\",\n  \"items\": 3,\n  \"order_id\": \"ORD-1234\",\n  \"status\": \"Processing\",\n  \"total_amount\": \"$111.19\"\n}"
    },
    {
      "type": "answer",
      "content": "The headphones with order ID ORD-1234 have been shipped via UPS with a tracking number TRK485533. The estimated delivery date is May 26, 2025, and the current status is \"Order Placed.\"\n\nThe order status is \"Processing,\" which means it has not been delivered yet. Since the order is still in processing, there is a possibility to change the delivery address. I recommend contacting customer support as soon as possible to request a change in the delivery address for your order."
    }
  ]
}
```
</details>
