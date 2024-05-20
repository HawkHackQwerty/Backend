from flask import Flask, request, jsonify
from flask_cors import CORS
from werkzeug.utils import secure_filename
from cohere_func import *
from resume_scraper import extract_text_from_pdf
from cover_scraper import extract_text_from_cover
from interview import interview_reply
import os

app = Flask(__name__)
CORS(app)

chat_history: list = []
job_title=""
job_desc=""

UPLOAD_FOLDER = 'uploads/'
os.makedirs(UPLOAD_FOLDER, exist_ok=True)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
ALLOWED_EXTENSIONS = {'pdf'}


@app.route("/jobs", methods=["GET", "POST"])
def get_jobs():
    global job_title, job_desc
    job_title = request.json.get('title')
    job_desc = request.json.get('description')
    
    return jsonify("Working!")

@app.route("/resume", methods=["GET", "POST"])
def resume():
    global job_desc
    if 'file' not in request.files:
        return jsonify({"error": "No file part"}), 400
    file = request.files['file']
    if file.filename == '':
        return jsonify({"error": "No selected file"}), 400
    if file:
        filename = secure_filename(file.filename)
        file_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
        file.save(file_path)
        resum = extract_text_from_pdf(file_path)
        feedback, score = get_resume_feedback(job_desc, resum, chat_history)
        return jsonify({"feedback": feedback, "score": score})
    else:
        return jsonify({"error": "Invalid file"}), 400

@app.route("/letter", methods=["GET", "POST"])
def letter():
    global job_desc, job_title
    if 'file' not in request.files:
        return jsonify({"error": "No file part"}), 400
    file = request.files['file']
    if file.filename == '':
        return jsonify({"error": "No selected file"}), 400
    if file:
        filename = secure_filename(file.filename)
        file_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
        file.save(file_path)
        cover = extract_text_from_cover(file_path)
        new_cv= get_new_cv(job_desc, job_title, cover, chat_history)
        return jsonify({"new_letter": new_cv})
    else:
        return jsonify({"error": "Invalid file"}), 400

@app.route("/video", methods=["GET", "POST"])
def video():
    if 'file' not in request.files:
        return jsonify({"error": "No file part"}), 400
    file = request.files['file']
    if file.filename == '':
        return jsonify({"error": "No selected file"}), 400
    if file:
        filename = secure_filename(file.filename)
        file_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
        file.save(file_path)

        reply = interview_reply(file_path, chat_history)

        return jsonify({"feedback": reply})
    else:
        return jsonify({"error": "Invalid file"}), 400


if __name__ == "__main__":
    app.run(debug=True)