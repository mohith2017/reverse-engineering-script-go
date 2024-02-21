import requests
import base64

application_email = "mohith.ny2024@gmail.com"
base_url = "https://ciphersprint.pulley.com/"
challenge_url = base_url + application_email
next_challenge_url = ""
# challenge_url = "https://ciphersprint.pulley.com/"


try:
    challenge = requests.get(challenge_url)
    
    # Check the status code
    if challenge.status_code == 200:
        response = challenge.json()
        print("Request successful!")
        print(challenge.text)
        next_challenge_url = response["encrypted_path"]

    else:
        print(f"Request failed with status code {challenge.status_code}")
except requests.exceptions.RequestException as e:
    print("Error making the request:", e)


print("Next challenge URL: ", next_challenge_url)
challenge_url = base_url + next_challenge_url
try:
    challenge = requests.get(challenge_url)
    
    # Check the status code
    if challenge.status_code == 200:
        response = challenge.json()
        print("Request successful!")
        print(challenge.text)
        next_challenge_url = response["encrypted_path"]

    else:
        print(f"Request failed with status code {challenge.status_code}")
except requests.exceptions.RequestException as e:
    print("Error making the request:", e)

print("Next challenge URL: ", next_challenge_url[5:])
decoded_path = base64.b64decode(next_challenge_url[5:]).decode('utf-8')
challenge_url = base_url + "task_" + decoded_path 
try:
    challenge = requests.get(challenge_url)
    
    # Check the status code
    if challenge.status_code == 200:
        response = challenge.json()
        print("Request successful!")
        print(challenge.text)
        next_challenge_url = response["encrypted_path"]

    else:
        print(f"Request failed with status code {challenge.status_code}")
except requests.exceptions.RequestException as e:
    print("Error making the request:", e)

print("Next challenge URL: ", next_challenge_url)