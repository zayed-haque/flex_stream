import os
import json
from quixstreams import Application, State
from quixstreams.models.serializers import JSONDeserializer, JSONSerializer
import google.generativeai as genai
from dotenv import load_dotenv
import psycopg2
from psycopg2.extras import Json

# Load environment variables
load_dotenv()

# Initialize the Quix Application
app = Application.Quix(
    broker_address="redpanda:29092",
    consumer_group="flex-stream-processing",
    auto_offset_reset="earliest",
)

# Define input and output topics
input_topic = app.topic("raw_data", value_deserializer=JSONDeserializer())
output_topic = app.topic("processed_data", value_serializer=JSONSerializer())

# Configure Gemini API
genai.configure(api_key=os.getenv("GOOGLE_API_KEY"))
model = genai.GenerativeModel("gemini-1.5-pro")

# Configure PostgreSQL connection
db_conn = psycopg2.connect(
    dbname=os.getenv("POSTGRES_DB"),
    user=os.getenv("POSTGRES_USER"),
    password=os.getenv("POSTGRES_PASSWORD"),
    host="postgres"
)


def store_processed_data(data):
    with db_conn.cursor() as cur:
        cur.execute("""
            INSERT INTO processed_data (original_id, data_type, original_data, processed_result, timestamp)
            VALUES (%s, %s, %s, %s, %s)
        """, (
            data['id'],
            data['data_type'],
            Json(data['original_data']),
            Json(data['processed_result']),
            data['timestamp']
        ))
    db_conn.commit()


def process_with_gemini(data, data_type):
    """
    Use Gemini 1.5 Flash to process data based on its type.
    """
    prompt = f"""
    Analyze the following {data_type} data and provide:
    1. A brief description or summary
    2. Key insights or patterns
    3. Any anomalies or unusual aspects
    4. Potential next steps or actions based on this data

    Data: {data}

    Response format:
    {{
        "description": "...",
        "insights": ["insight1", "insight2", ...],
        "anomalies": ["anomaly1", "anomaly2", ...],
        "next_steps": ["step1", "step2", ...]
    }}
    """

    response = model.generate_content(prompt)

    try:
        result = json.loads(response.text)
    except json.JSONDecodeError:
        # If the response is not valid JSON, return a default structure
        result = {
            "description": "Failed to generate analysis",
            "insights": [],
            "anomalies": [],
            "next_steps": [],
        }

    return result


def process_data(value: dict, state: State):
    """
    Process the incoming data using Gemini 1.5 Flash for analysis.
    """
    data_type = value.get("data_type", "unknown")
    data = value.get("data")

    if data is not None:
        # Process the data using Gemini
        processed_result = process_with_gemini(data, data_type)

        # Create processed data dictionary
        processed_data = {
            "id": value.get("id", ""),
            "original_data": data,
            "data_type": data_type,
            "processed_result": processed_result,
            "timestamp": value.get("timestamp", ""),
        }

        # Keep track of processed count
        current_count = state.get("processed_count", 0)
        state["processed_count"] = current_count + 1
        processed_data["processed_count"] = state["processed_count"]

        store_processed_data(processed_data)

        return processed_data
    else:
        # If there's no data field, return the original value
        return value


# Set up the stream processing pipeline
(input_topic.process(process_data).to(output_topic))

if __name__ == "__main__":
    app.run()
