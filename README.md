# Yordamchi

#### Video Demo: https://youtu.be/PJPSVEC773E

In Uzbekistan, there is a problem with using ChatGPT: OpenAI's service is not available in our country. I wanted to
provide a solution for this problem by building a Telegram bot that is connected to OpenAI API. Why Telegram? Because it
is the most popular messaging app in Uzbekistan.

At the time of writing, the bot has more than 13 000 users and there are many who write grateful messages to me,
mentioning it is of great help. I am happy that I could help people with this project. Actually, "yordamchi" means "
assistant" or "helper" in Uzbek.

You can try yourself: https://t.me/yordamchi_ai_bot

The project indeed began as the final project for CS50x course. But it grew into something big and unforeseen.

## Components

Project started small and simple, but over time new features were added, and it kept gaining new components. Here is how
the project is overall structured:

- **redis**: Redis is used as the quick storage means. It stores users' status, context of the conversations, and other
  data that is needed for the bot to function. It is also used to record the current activity making it possible to
  restore the activity in case of a crash or restart.
- **lib**: This is a library of constants, functions and types that are used in the project. It also contains variables
  for handling errors.
- **payme**: This is a payment system that is used to accept payments for the service. It is not used at the moment of
  writing.
- **client**: These are the clients that are used to connect to OpenAI API, Telegram API, Azure Computer Vision API,
  translator and tokenizer. It is responsible for sending requests and
  receiving responses.
- **postgres**: This is the database that is used to store users' data. It is used to store users' and messages'
  metadata. The messages themselves are not stored and not shared with anyone.
- **text**: This is a library that contains text of messages that are sent to users. The bot supports three languages:
  Uzbek, Russian and English.
- **processor**: This is the core of the project. It contains the main logic of the bot. It is responsible for handling
  messages, processing them, and sending responses.

## Contact Information

- **Email**: doniyorbek752@gmail.com
- **Telegram**: https://t.me/dro_14
- **LinkedIn**: https://linkedin.com/in/dro14
- **Github**: https://github.com/dro14


