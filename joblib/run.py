import joblib

model_file = "model.pk1"

model_dict = joblib.load(model_file)
vectorizer = model_dict['vectorizer']
classifier = model_dict['classifier']

user_input = input("Enter a Linux command or description: ").strip()

X_input = vectorizer.transform([user_input])

predicted_label = classifier.predict(X_input)[0]

print("Predicted category:", predicted_label)
