import asyncio
import zmq.asyncio
import struct
from dataclasses import dataclass

@dataclass
class ReturnMessage:
    ResumeFeedback: str
    Score: int
    UpdatedCoverLetter: str
    LinkedInColdMessage: str

async def CreateServerSocket():
    """
    Creates and binds an asynchronous REP socket for the server
    """
    context = zmq.asyncio.Context()
    socket = context.socket(zmq.REP)
    socket.bind("tcp://*:5555")
    return socket, context

async def CloseServerSocket(socket, context):
    """Closes the socket and terminates the ZeroMQ context asynchronously."""
    socket.close()
    context.term()


async def ReceiveJobInformation(socket):
    """
    Receives binary data asynchronously and unpacks it into job information
    """
    binary_data = await socket.recv()
    # Read job title length and job title
    job_title_length = struct.unpack('<I', binary_data[:4])[0]
    job_title = binary_data[4:4 + job_title_length].decode('utf-8')

    # Read job description length and job description
    job_description_start = 4 + job_title_length
    job_description_length = struct.unpack('<I', binary_data[job_description_start:job_description_start + 4])[0]
    job_description = binary_data[job_description_start + 4:job_description_start + 4 + job_description_length].decode('utf-8')

    return job_title, job_description

async def ProcessRequest(job_title, job_description):
    """
    Processes job information and prepares feedback asynchronously
    """

    # PLACEHOLDER FOR ACTUAL PROCESSING
    # YOUR FUNCTION SHOULD RETURN A ReturnMessage OBJECT
    # example:
    # feedback_data = func() which returns a ReturnMessage object

    # Below is just for testing
    feedback_data = ReturnMessage(
        ResumeFeedback="Looks great",
        Score=85,
        UpdatedCoverLetter="Updated cover letter based on job description",
        LinkedInColdMessage="Hi, I noticed we share similar interests..."
    )
    return feedback_data

async def SendFeedbackPackage(socket, feedback_data):
    """
    Serializes feedback into binary and sends it back to the client asynchronously
    """
    resume_feedback = feedback_data.ResumeFeedback.encode('utf-8')
    updated_cover_letter = feedback_data.UpdatedCoverLetter.encode('utf-8')
    cold_message = feedback_data.LinkedInColdMessage.encode('utf-8')
    score = feedback_data.Score
    reply_data = struct.pack(f'<{len(resume_feedback)}sB{len(updated_cover_letter)}s{len(cold_message)}s', 
                             resume_feedback, score, updated_cover_letter, cold_message)
    await socket.send(reply_data)

async def main():
    socket, context = await CreateServerSocket()
    try:
        while True:
            job_title, job_description = await ReceiveJobInformation(socket)
            feedback_data = await ProcessRequest(job_title, job_description)
            await SendFeedbackPackage(socket, feedback_data)
    except KeyboardInterrupt:
        print("Shutting down server...")
    finally:
        await CloseServerSocket(socket, context)

if __name__ == "__main__":
    asyncio.run(main())
