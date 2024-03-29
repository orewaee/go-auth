## Attention

This app was written for educational purposes only and may contain errors.
Ideas and development methods have also been applied to explore them.
The work logic and goals of the project are subject to change during its existence.


## About

The application is an authorization based on access and refresh [JWT](https://github.com/golang-jwt/jwt) tokens.
It requires an email confirmation of the action to create a user.
All data is stored in [MongoDB](https://github.com/mongodb/mongo). [Fiber](https://github.com/gofiber/fiber) is chosen as the web framework.


## Pull Requests or Commits

You are encouraged to use the following prefixes for Pull Requests or Commits:

> 🔥 feat, 🩹 fix, 📚 docs, 🎨 style, ♻️ refactor, 🚀 perf, 🚨 test, 🔨 chore


## Routes

#### Public

- `GET /ping`
- `POST /signup`
- `POST /activate/:secret`
- `POST /signin`
- `POST /refresh`

#### Protected

- `GET /user`
