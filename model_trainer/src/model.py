from sklearn.ensemble import RandomForestClassifier
from sklearn.feature_extraction.text import TfidfVectorizer


def create_model():
    return RandomForestClassifier(n_estimators=100, random_state=42)


def preprocess_data(df):
    # Create TF-IDF features from the description
    vectorizer = TfidfVectorizer(max_features=1000)
    X = vectorizer.fit_transform(df['description'])
    y = df['sentiment']
    return X, y, vectorizer
