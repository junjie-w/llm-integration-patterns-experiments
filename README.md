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
