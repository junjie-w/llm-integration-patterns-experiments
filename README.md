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

### 5. Multi Agent

**Endpoint**: `POST /api/support/multi-agent`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "message": "I ordered wireless headphones last week but they won't connect to my laptop. The order number is ORD-1234."
}
```

**Example Response**
```json
{
    "conversation_id": "ce6308c5-47dc-426b-aa26-b72455924555",
    "reply": "Hello,\n\nI'm sorry to hear that you are facing difficulty connecting your wireless headphone to your laptop. To troubleshoot this issue, please follow these steps provided by TechnicalSupport:\n\n1. Ensure the wireless headphone is fully charged.\n2. Put the headphone in pairing mode by following the instructions in the user manual.\n3. On your laptop, navigate to the Bluetooth settings and turn on Bluetooth.\n4. Locate the wireless headphone in the list of available devices and select it to pair.\n5. Enter the default PIN code (usually 0000 or 1234) if required.\n6. Once paired successfully, your laptop should show a notification confirming the connection.\n7. Test the audio output to verify if the headphone is working correctly.\n\nIf you encounter any difficulties, please contact the technical support team at techsupport@example.com for further assistance.\n\nFor updates on your order with the number ORD-1234, I recommend reaching out to our customer support team for the latest information.\n\nShould you have any other questions or concerns, please feel free to let me know.",
    "agents": [
        "TechnicalSupport",
        "CustomerSupport"
    ],
    "complete": true
}
```
</details>

### 6. Evaluation

**Evaluation Endpoint**: `POST /api/support/evaluate`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```json
{
  "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
  "pattern_types": ["basic_llm_completion", "knowledge_rag", "function_calling", "reasoning_agent", "multi_agent"]
}
```

**Example Response**
```json
{
    "report_id": "786ec6ab-0845-48a2-8aaf-c4e96834ea2e",
    "report": {
        "id": "786ec6ab-0845-48a2-8aaf-c4e96834ea2e",
        "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
        "timestamp": "2025-05-25T13:56:08.889893+02:00",
        "results": [
            {
                "id": "4f4652ef-4334-42ba-b3a6-ac866661a83e",
                "pattern_type": "basic_llm_completion",
                "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
                "response": "I apologize for the delay in receiving your headphones. Let me check the status of your order with order number ORD-1234 and get back to you shortly. Thank you for your patience.",
                "response_time_ms": 694,
                "auto_rating": 0.7,
                "evaluation_notes": "\nExplanation: The response is relevant as it acknowledges the delay in receiving the headphones and mentions checking the status of the order number provided. It shows accuracy by referring to the specific order number. While it is not yet complete in providing a solution, it assures the customer that they will get back shortly. The response is clear and offers reassurance, which is helpful in managing customer expectations."
            },
            {
                "id": "bb989249-7854-4cb2-b219-5038c02cb729",
                "pattern_type": "knowledge_rag",
                "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
                "response": "I'm sorry to hear about the delay in receiving your headphones. To investigate further, I recommend contacting our customer support team with your order number ORD-1234. They will assist you in tracking your order and providing an update on its status.",
                "response_time_ms": 1161,
                "auto_rating": 0.8,
                "evaluation_notes": "\nExplanation: The response is relevant as it directly acknowledges the delay in receiving the headphones and provides a solution to track the order with the order number provided. The information appears accurate, complete, and clear, guiding the customer on the next steps to take. It is helpful by offering a specific course of action to address the issue. The response could be improved by including an apology for the inconvenience caused."
            },
            {
                "id": "d550e8f0-8ba8-49bc-9206-2c75e2daf0c4",
                "pattern_type": "function_calling",
                "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
                "response": "I found the information regarding your order:\n\n- Order Status: **Cancelled**\n- Shipping Carrier: **FedEx**\n- Tracking Number: **TRK834061**\n- Estimated Delivery Date: **2025-05-28**\n\nIt seems like your order has been cancelled. If you have any further questions or need assistance, feel free to let me know.",
                "response_time_ms": 1461,
                "auto_rating": 0.6,
                "evaluation_notes": "\nExplanation: The response partially addresses the query by providing relevant information about the order status being cancelled. However, it lacks accuracy as it states the estimated delivery date in the future instead of acknowledging the cancellation. It is somewhat complete but could be clearer by explicitly stating the reason for cancellation and offering assistance on next steps. Overall, it falls in the average category for addressing the query with a few issues."
            },
            {
                "id": "a586cd14-5dd7-4153-bb2a-1b5125d24b70",
                "pattern_type": "reasoning_agent",
                "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
                "response": "The headphones you ordered with the order number ORD-1234 have been shipped via FedEx. The current status is \"Out for Delivery,\" and the estimated delivery date is May 28, 2025. You can track your package using the tracking number TRK411121 provided.",
                "response_time_ms": 1267,
                "auto_rating": 0.9,
                "evaluation_notes": "\nExplanation: The response is excellent. It directly addresses the query by providing the current status of the order, the shipping information, and ways to track the package. The information is accurate, complete, clear, and highly helpful to the customer."
            },
            {
                "id": "31b4193a-dd4a-4151-b1bd-7a971fcc83b6",
                "pattern_type": "multi_agent",
                "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
                "response": "I apologize for the delay in receiving your headphones. Let me check the status of your order ORD-1234 for you. Please allow me a moment to investigate this for you.",
                "response_time_ms": 503,
                "auto_rating": 0.9,
                "evaluation_notes": "\nExplanation: The response addresses the query directly by acknowledging the delay in receiving the headphones and mentioning the order number. The customer support representative also assures the customer of checking the status of the order, which shows readiness to assist. The response is clear, concise, and likely to be helpful in resolving the customer's issue promptly."
            }
        ]
    }
}
```
</details>


**Report Retrieval Endpoint**: `GET /api/support/evaluate/report/:id`

<details>
<summary><strong>Example Request/Response</strong></summary>

**Example Request**
```
GET /api/support/evaluate/report/786ec6ab-0845-48a2-8aaf-c4e96834ea2e
```

**Example Response**
```json
{
    "id": "786ec6ab-0845-48a2-8aaf-c4e96834ea2e",
    "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
    "timestamp": "2025-05-25T13:56:08.889893+02:00",
    "results": [
        {
            "id": "4f4652ef-4334-42ba-b3a6-ac866661a83e",
            "pattern_type": "basic_llm_completion",
            "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
            "response": "I apologize for the delay in receiving your headphones. Let me check the status of your order with order number ORD-1234 and get back to you shortly. Thank you for your patience.",
            "response_time_ms": 694,
            "auto_rating": 0.7,
            "evaluation_notes": "\nExplanation: The response is relevant as it acknowledges the delay in receiving the headphones and mentions checking the status of the order number provided. It shows accuracy by referring to the specific order number. While it is not yet complete in providing a solution, it assures the customer that they will get back shortly. The response is clear and offers reassurance, which is helpful in managing customer expectations."
        },
        {
            "id": "bb989249-7854-4cb2-b219-5038c02cb729",
            "pattern_type": "knowledge_rag",
            "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
            "response": "I'm sorry to hear about the delay in receiving your headphones. To investigate further, I recommend contacting our customer support team with your order number ORD-1234. They will assist you in tracking your order and providing an update on its status.",
            "response_time_ms": 1161,
            "auto_rating": 0.8,
            "evaluation_notes": "\nExplanation: The response is relevant as it directly acknowledges the delay in receiving the headphones and provides a solution to track the order with the order number provided. The information appears accurate, complete, and clear, guiding the customer on the next steps to take. It is helpful by offering a specific course of action to address the issue. The response could be improved by including an apology for the inconvenience caused."
        },
        {
            "id": "d550e8f0-8ba8-49bc-9206-2c75e2daf0c4",
            "pattern_type": "function_calling",
            "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
            "response": "I found the information regarding your order:\n\n- Order Status: **Cancelled**\n- Shipping Carrier: **FedEx**\n- Tracking Number: **TRK834061**\n- Estimated Delivery Date: **2025-05-28**\n\nIt seems like your order has been cancelled. If you have any further questions or need assistance, feel free to let me know.",
            "response_time_ms": 1461,
            "auto_rating": 0.6,
            "evaluation_notes": "\nExplanation: The response partially addresses the query by providing relevant information about the order status being cancelled. However, it lacks accuracy as it states the estimated delivery date in the future instead of acknowledging the cancellation. It is somewhat complete but could be clearer by explicitly stating the reason for cancellation and offering assistance on next steps. Overall, it falls in the average category for addressing the query with a few issues."
        },
        {
            "id": "a586cd14-5dd7-4153-bb2a-1b5125d24b70",
            "pattern_type": "reasoning_agent",
            "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
            "response": "The headphones you ordered with the order number ORD-1234 have been shipped via FedEx. The current status is \"Out for Delivery,\" and the estimated delivery date is May 28, 2025. You can track your package using the tracking number TRK411121 provided.",
            "response_time_ms": 1267,
            "auto_rating": 0.9,
            "evaluation_notes": "\nExplanation: The response is excellent. It directly addresses the query by providing the current status of the order, the shipping information, and ways to track the package. The information is accurate, complete, clear, and highly helpful to the customer."
        },
        {
            "id": "31b4193a-dd4a-4151-b1bd-7a971fcc83b6",
            "pattern_type": "multi_agent",
            "query": "I ordered headphones last week but they haven't arrived yet. My order number is ORD-1234.",
            "response": "I apologize for the delay in receiving your headphones. Let me check the status of your order ORD-1234 for you. Please allow me a moment to investigate this for you.",
            "response_time_ms": 503,
            "auto_rating": 0.9,
            "evaluation_notes": "\nExplanation: The response addresses the query directly by acknowledging the delay in receiving the headphones and mentioning the order number. The customer support representative also assures the customer of checking the status of the order, which shows readiness to assist. The response is clear, concise, and likely to be helpful in resolving the customer's issue promptly."
        }
    ]
}
```
</details>
