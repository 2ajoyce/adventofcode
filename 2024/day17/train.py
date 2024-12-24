import keras.api
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
import tensorflow as tf
from tensorflow.python.client import device_lib
import keras as keras
from keras import layers
from keras import utils
from keras import callbacks
from keras import backend
from keras import models
import matplotlib.pyplot as plt

# Constants
MAX_SEQUENCE_LENGTH = 25
NUM_CLASSES = 10  # Digits 0-9
EMBEDDING_DIM = 32
BATCH_SIZE = 128
EPOCHS = 50
MAX_INT64 = 9223372036854775807  # Maximum value for int64


def encode_output(sequence, max_length=MAX_SEQUENCE_LENGTH):
    """
    Encodes a digit sequence string into a one-hot encoded NumPy array.

    Parameters:
        sequence (str): The digit sequence as a string.
        max_length (int): The fixed length to pad the sequence.

    Returns:
        np.ndarray: One-hot encoded array of shape (max_length, NUM_CLASSES)
    """
    sequence = sequence.rjust(max_length, "0")  # Left-padding
    # Alternatively, use right-padding
    # sequence = sequence.ljust(max_length, '0')  # Right-padding
    encoded = keras.utils.to_categorical(
        [int(char) for char in sequence], num_classes=NUM_CLASSES
    )
    return encoded


def decode_sequence(encoded_seq):
    """
    Decodes a one-hot encoded sequence back to a digit string.

    Parameters:
        encoded_seq (np.ndarray): One-hot encoded array of shape (max_length, NUM_CLASSES)

    Returns:
        str: Decoded digit sequence as a string.
    """
    digits = np.argmax(encoded_seq, axis=-1)
    # Convert to string and remove leading zeros
    digit_str = "".join(str(digit) for digit in digits).lstrip("0")
    return digit_str if digit_str else "0"


def load_and_prepare_data(csv_file, sample_size=None, random_state=42):
    """
    Loads data from a CSV file and encodes input and output sequences.

    Parameters:
        csv_file (str): Path to the CSV file.
        sample_size (int, optional): Number of rows to sample. If None, load all data.
        random_state (int, optional): Seed for reproducibility.

    Returns:
        tuple: Encoded output sequences (X) and normalized input scalars (y).
    """
    # Load the dataset with input and output as strings
    data = pd.read_csv(csv_file, dtype={"input": str, "output": str})

    # If a sample size is specified, randomly sample that many rows
    if sample_size is not None:
        if sample_size > len(data):
            raise ValueError(
                f"Sample size {sample_size} exceeds total data size {len(data)}."
            )
        data = data.sample(n=sample_size, random_state=random_state).reset_index(
            drop=True
        )
        print(f"Sampled {sample_size} rows from the dataset.")

    # Encode the output sequences (already correct)
    X = np.array([encode_output(seq) for seq in data["output"].values])

    # Encode the input integers as scalar values
    y = data["input"].astype(np.int64).values

    # Normalize y to be between 0 and 1
    y = y / MAX_INT64

    return X, y


def build_model():
    """
    Builds and compiles the model to map scalar inputs to output sequences.

    Returns:
        tensorflow.keras.Model: Compiled Keras model.
    """
    model = keras.models.Sequential()

    # Input layer for scalar input
    model.add(layers.Input(shape=(1,)))  # Scalar input

    # Dense layers to process scalar input
    model.add(layers.Dense(128, activation="relu"))
    model.add(layers.Dense(256, activation="relu"))
    model.add(layers.Dense(MAX_SEQUENCE_LENGTH * NUM_CLASSES, activation="relu"))

    # Reshape to match output sequence shape
    model.add(layers.Reshape((MAX_SEQUENCE_LENGTH, NUM_CLASSES)))

    # LSTM layer to process the reshaped input
    model.add(layers.LSTM(128, return_sequences=True))
    model.add(layers.Dropout(0.2))

    # TimeDistributed Dense layer with softmax activation for output
    model.add(layers.TimeDistributed(layers.Dense(NUM_CLASSES, activation="softmax")))

    # Compile the model
    loss = keras.losses.CategoricalCrossentropy()
    optimizer = keras.optimizers.Adam()
    catMetric = keras.metrics.CategoricalAccuracy()
    model.compile(
        loss=loss,
        optimizer=optimizer,
        metrics=[catMetric],
    )
    return model


def plot_history(history):
    """
    Plots the training and validation loss and accuracy.

    Parameters:
        history (History): Keras History object.
    """
    # Plot loss
    plt.figure(figsize=(12, 4))
    plt.subplot(1, 2, 1)
    plt.plot(history.history["loss"], label="Train Loss")
    plt.plot(history.history["val_loss"], label="Validation Loss")
    plt.title("Loss Over Epochs")
    plt.xlabel("Epoch")
    plt.ylabel("Loss")
    plt.legend()

    # Plot accuracy
    plt.subplot(1, 2, 2)
    plt.plot(history.history["categorical_accuracy"], label="Train Accuracy")
    plt.plot(history.history["val_categorical_accuracy"], label="Validation Accuracy")
    plt.title("Accuracy Over Epochs")
    plt.xlabel("Epoch")
    plt.ylabel("Accuracy")
    plt.legend()

    plt.tight_layout()
    plt.show()


def main():
    print("Num GPUs Available: ", len(tf.config.list_physical_devices("GPU")))
    print(device_lib.list_local_devices())

    # Define the sample size
    SAMPLE_SIZE = None  # Set to None to use the full dataset

    # Load and prepare data with sampling
    X, y = load_and_prepare_data("data.csv", sample_size=SAMPLE_SIZE)
    print(
        f"Output (X) shape: {X.shape}"
    )  # (sample_size, MAX_SEQUENCE_LENGTH, NUM_CLASSES)
    print(f"Input (y) shape: {y.shape}")  # (sample_size,)

    # Split the data
    X_train, X_temp, y_train, y_temp = train_test_split(
        X, y, test_size=0.3, random_state=42
    )
    X_val, X_test, y_val, y_test = train_test_split(
        X_temp, y_temp, test_size=0.5, random_state=42
    )

    print(f"Training samples: {X_train.shape[0]}")
    print(f"Validation samples: {X_val.shape[0]}")
    print(f"Test samples: {X_test.shape[0]}")

    # Build the model
    model = build_model()
    model.summary()

    # Define callbacks
    early_stopping = keras.callbacks.EarlyStopping(
        monitor="val_loss", patience=5, restore_best_weights=True
    )
    model_checkpoint = keras.callbacks.ModelCheckpoint(
        "input_predictor_model.keras", save_best_only=True, monitor="val_loss"
    )

    # Train the model
    history = model.fit(
        y_train,  # Input: scalar values
        X_train,  # Output: sequences
        validation_data=(y_val, X_val),
        epochs=EPOCHS,
        batch_size=BATCH_SIZE,
        callbacks=[early_stopping, model_checkpoint],
    )

    # Plot training history
    plot_history(history)

    # Evaluate the model
    test_loss, test_accuracy = model.evaluate(y_test, X_test, verbose=0)
    print(f"Test Loss: {test_loss:.4f}")
    print(f"Test Accuracy: {test_accuracy:.4f}")

    # Save the model
    model.save("best_model.keras")

    # Example prediction
    sample_input = y_test[0]  # Normalized scalar input
    actual_output = X_test[0]  # Encoded output sequence

    # Expand dimensions to match model's expected input shape
    sample_input_expanded = np.expand_dims(sample_input, axis=0)  # Shape: (1,)

    # Predict the output sequence from the input
    predicted_output = model.predict(sample_input_expanded)

    # Decode the predicted sequence
    predicted_sequence = decode_sequence(predicted_output[0])

    # Decode the actual output sequence
    actual_sequence = decode_sequence(actual_output)

    # Convert the normalized input back to the original integer
    actual_input_integer = int(sample_input * MAX_INT64)

    print("\n--- Example Prediction ---")
    print(f"Input Integer: {actual_input_integer}")
    print(f"Actual Output Sequence: {actual_sequence}")
    print(f"Predicted Output Sequence: {predicted_sequence}")


if __name__ == "__main__":
    main()
