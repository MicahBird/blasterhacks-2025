import socket
import joblib
import os, os.path

def run_model(command):
    print(command)
    model_file = "model.pk1"

    model_dict = joblib.load(model_file)
    vectorizer = model_dict['vectorizer']
    classifier = model_dict['classifier']

    user_input = command.strip()

    X_input = vectorizer.transform([user_input])

    predicted_label = classifier.predict(X_input)[0]
    
    print(predicted_label)
    return predicted_label


def receive():
    if os.path.exists("/tmp/r.sock"):
      os.remove("/tmp/r.sock")

    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind("/tmp/r.sock")
    while True:
        server.listen(1)
        conn, addr = server.accept()
        conn.setblocking(True)
        output = str(run_model(str(conn.recvmsg(2048))))
        send(output)


def send(label):
    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.connect("/tmp/s.sock")
    print("sending tag back")
    server.send(bytes(label, encoding='utf8'))

# main but I don't feel like doing the name == main thing
print("Starting server")
receive()
