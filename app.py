from flask import Flask, render_template, request, jsonify
from flask_cors import CORS
import language_tool_python

app = Flask(__name__)
CORS(app)  # Allow cross-origin requests

tool = language_tool_python.LanguageTool('en-US')

# Render index.html when accessing the root URL
@app.route('/')
def index():
    return render_template('index.html')

@app.route('/correct', methods=['POST'])
def correct_grammar():
    data = request.get_json()
    text = data.get('text', '')

    if not text:
        return jsonify({'error': 'No text provided'}), 400
    
    matches = tool.check(text)
    corrected_text = language_tool_python.utils.correct(text, matches)
    
    return jsonify({'corrected_text': corrected_text})

if __name__ == '__main__':
    app.run(debug=True)
