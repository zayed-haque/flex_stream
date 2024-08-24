import schedule
import time
from data_loader import load_data
from model import create_model, preprocess_data
from model_trainer.src.trainer import train_model


def training_job():
    print("Starting training job...")
    data = load_data()
    preprocessed_data = preprocess_data(data)
    model = create_model()
    train_model(model, preprocessed_data)
    print("Training job completed.")


schedule.every().day.at("02:00").do(training_job)

if __name__ == "__main__":
    print("Model trainer service started.")
    training_job()
    while True:
        schedule.run_pending()
        time.sleep(60)
