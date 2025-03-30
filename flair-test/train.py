import pandas as pd
import flair
from flair.data import Corpus
from flair.data import Sentence
from flair.models import SequenceTaggerModel
from flair.datasets import ColumnCorpus

# Load the CSV file
def load_data(csv_file):
    """Load the CSV file containing the Linux commands and tags"""
    data = pd.read_csv(csv_file)
    return data

# Prepare the data for training
def prepare_data(data):
    """Prepare the data for training by converting it to Flair format"""
    sentences = []
    for index, row in data.iterrows():
        sentence = Sentence(row['command'])
        sentence.add_tag('ner', row['tag'])
        sentences.append(sentence)
    return sentences

# Create a Corpus object
def create_corpus(sentences):
    """Create a Corpus object from the prepared sentences"""
    corpus = ColumnCorpus(None, ColumnCorpus._get_sentences, sentences=sentences)
    return corpus

# Train the model
def train_model(corpus):
    """Train a SequenceTaggerModel on the Corpus object"""
    tagger = SequenceTaggerModel('ner')
    tagger.train(corpus, learning_rate=0.1, mini_batch_size=32, max_epochs=10)
    return tagger

# Main function
def main():
    csv_file = 'commands.csv'  # Replace with your CSV file
    data = load_data(csv_file)
    sentences = prepare_data(data)
    corpus = create_corpus(sentences)
    tagger = train_model(corpus)
    # Save the trained model
    tagger.save('linux_command_tagger')

if __name__ == '__main__':
    main()
