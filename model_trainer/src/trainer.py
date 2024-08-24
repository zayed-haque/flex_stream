import joblib
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report


def train_model(model, data):
    X, y, vectorizer = data
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    model.fit(X_train, y_train)
    y_pred = model.predict(X_test)
    print(classification_report(y_test, y_pred))

    joblib.dump(model, '/app/models/trained_model.joblib')
    joblib.dump(vectorizer, '/app/models/vectorizer.joblib')

    print("Model and vectorizer saved successfully.")
