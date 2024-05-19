import zmq
import threading
from io import BytesIO
import json

from cohere_func import *

def handler_resume(message, port):
    try:
        data = json.loads(message)
        user_id = data['userID']
        file_data = BytesIO(data['fileData'])
        print(f"Received resume file from user {user_id} on port {port}")

        # # Use PyPDF2 to read the PDF from BytesIO
        # pdf_reader = PdfReader(file_data)
        # number_of_pages = len(pdf_reader.pages)
        # print(f"Received resume file from user {user_id} on port {port}, Number of pages: {number_of_pages}")

        # Example processing: Just read the text from the first page
        # first_page_text = pdf_reader.pages[0].extract_text()
        # response_part1 = f"Resume processed for user {user_id}, first page text: {first_page_text[:100]}..

        # Process resume
        response_part1 = f"Resume processed for user {user_id}"
        response_part2 = "Verification Complete"
        response_part3 = "Ready for next step"

    except json.JSONDecodeError as e:
        print(f"Error decoding JSON on port {port}: {e}")
        response_part1, response_part2, response_part3 = "Failed to decode JSON", "", ""
    except KeyError as e:
        print(f"Missing key in data on port {port}: {e}")
        response_part1, response_part2, response_part3 = "Data is missing keys", "", ""

    # Return three strings as a multipart response
    return [response_part1, response_part2, response_part3]

def handler_cover(message, port):
    try:
        data = json.loads(message)
        user_id = data['userID']
        file_data = BytesIO(data['fileData'])
        print(f"Received file from user {user_id} on port {port}")

        # Placeholder for processing the PDF
        response_message = "PDF processing completed successfully"
        return [response_message]

    except json.JSONDecodeError as e:
        print(f"Error decoding JSON on port {port}: {e}")
        return ["Failed to decode JSON"]
    except KeyError as e:
        print(f"Missing key in data on port {port}: {e}")
        return ["Data is missing keys"]

def handler_video(message, port):
    try:
        data = json.loads(message)
        user_id = data['userID']
        video_data = BytesIO(data['fileData'])   
        video_filename = f"received_video_port_{port}.mp4"
        with open(video_filename, 'wb') as video_file:
            video_file.write(video_data.getvalue())
        print(f"Saved video to {video_filename}, size {len(video_data.getvalue())} bytes")
        return ["Video saved successfully"]
    except Exception as e:
        print(f"Error processing video on port {port}: {e}")
        return ["Failed to process video"]

def handler_job(message, port):
    try:
        data = json.loads(message)
        job_title = data['stringOne']
        job_desc = data['stringTwo']
        user_id = data['userID']
        print(f"Received on port {port}: stringOne={job_title}, stringTwo={job_desc}")

        response_message = "Job info processed successfully"
        return [response_message]

    except json.JSONDecodeError as e:
        print(f"Error decoding JSON on port {port}: {e}")
        return ["Failed to decode JSON"]
    except KeyError as e:
        print(f"Missing key in JSON data on port {port}: {e}")
        return ["Data is missing keys"]

def rep_socket(port, handler):
    context = zmq.Context()
    socket = context.socket(zmq.REP)
    socket.bind(f"tcp://*:{port}")

    while True:
        message = socket.recv_string()
        response_parts = handler(message, port)
        socket.send_multipart([part.encode() for part in response_parts])

if __name__ == "__main__":
    port_handlers = {
        "5550": handler_resume,
        "5551": handler_cover,
        "5552": handler_video,
        "5553": handler_job,
    }

    threads = []
    for port, handler in port_handlers.items():
        thread = threading.Thread(target=rep_socket, args=(port, handler))
        thread.start()
        threads.append(thread)

    for thread in threads:
        thread.join() 
