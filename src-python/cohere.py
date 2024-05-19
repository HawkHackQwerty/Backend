#from parser import letter, resume
#from omar import job dest, and title 
from resume import fahmi_resume, krish_resume
import cohere 

#put in .env file meow
co = cohere.Client("client")

def get_reply(prompt: str, history: list) -> str:

    response = co.chat(message=prompt, chat_history=history)

    answer = response.text

    history.append({"user_name": "User", "text": prompt})
    history.append({"user_name": "Chatbot", "text": answer})

    return answer

def preset_for_resume(history: list):
    
    get_reply("Here is a good example of a resume, take it as a note, you are going to be a resume expert who is reviewing a resume we are going to give you, here is the first resume: {}".format(fahmi_resume), history)
    
    get_reply("Here is another great example of a resume, note that the main great thing they did is that they included measurable metrics into their job descriptions {}".format(krish_resume), history)

    get_reply("Make sure to take some notes on what makes a good resume, for example: having separate sections, including measured metrics to demonstrate how it went, and being very concise, I will now give you a resume, take it, and give feedback to it in a bullet point format, and give it a score out of ten at the end and separate the feedback and score with a |", history)


def get_resume_feedback(job_desc, resume, history):
    preset_for_resume(history)
    return get_reply("here is the resume you are reviewing,: {}, and here is the job it is for {}".format(resume, job_desc)).split('|')

def get_new_cv(job_desc, job_title, cv, history):
    return get_reply("Thanks, now take that information from their resume, use this job description: {}, and this job title {}, and edit this cover letter that is from the user: {}. ONLY CHANGE WHAT IS INSIDE THE CURLY BRACKETS { }, EVERYTHING ELSE THEY WANT TO KEEP THE SAME!! ONLY RETURN THE ACTUAL LETTER, DO NOT ADD ON ANY EXTRA DESCRIPTION OR ANYTHING".format(job_desc, job_title, cv), history)

def get_cold_email(history):
    return get_reply("Now that you know all about the person, write a little template for a short and concise (300 char max) cold email they can send to a recruiter that they found, make sure to take into account the job description, as well as the persons personal strengths, ONLY GIVE BACK THE TEMPLATE, DO NOT GIVE ANYTHING ELSE", history=history)
