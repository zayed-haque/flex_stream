import pandas as pd
from sqlalchemy import create_engine
import os


def load_data():
    db_url = f"postgresql://{os.environ['POSTGRES_USER']}:{os.environ['POSTGRES_PASSWORD']}@postgres/{os.environ['POSTGRES_DB']}"
    engine = create_engine(db_url)

    query = """
    SELECT original_id, data_type,
           processed_result->>'description' as description,
           processed_result->>'sentiment' as sentiment,
           timestamp
    FROM processed_data
    WHERE timestamp > NOW() - INTERVAL '7 days'
    """
    df = pd.read_sql(query, engine)

    df['sentiment'] = df['sentiment'].map({'positive': 1, 'negative': 0, 'neutral': 0.5})

    return df
