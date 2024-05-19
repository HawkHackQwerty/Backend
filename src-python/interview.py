from gemeni import get_reply
import moviepy.editor as mp
import speech_recognition as sr


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


# get the video from frontend
clip = mp.VideoFileClip(video_name)

audio = clip.audio

audio.write_audiofile("interview.wav")

interview_audio_text = audio_to_text("interview.wav")

# cohere call, get the 
history = []
interview_feedback = get_reply(prompt="Here is someone who did a video interview question and answered it, I transcribed the text for you, give them feedback on how they can answer the questions better: Their questionw was: {}, and their response was: {}, give them that feedback!".format(question, interview_audio_text), history=history)
