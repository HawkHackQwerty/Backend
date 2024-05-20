from cohere_func import get_reply
import speech_recognition as sr
import ffmpeg


def audio_to_text(filename):
    recognizer = sr.Recognizer()

    text = sr.AudioFile(filename)
    with text as source:
        audio = recognizer.record(source)
    try:
        speech = recognizer.recognize_google(audio)
        return speech
    except Exception as e:
        print("Exception: "+str(e))

def interview_reply(filename, history):
    (
    ffmpeg.input(filename)
    .output("interview.wav")
    .run()
    )
    
    interview_audio_text = audio_to_text("interview.wav")
    return get_reply(prompt="Here is someone who did a video interview question and answered it, I transcribed the text for you, give them feedback on how they can answer the questions better: their question was Tell me about yourself. and their response was: {}, give them that feedback! PLEASE LIMIT YOUR RESPONSE TO 200 WORDS MAX!".format(interview_audio_text), history=history)