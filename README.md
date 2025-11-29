# Basic authorization system created in go
This project is meant to be a learning experience for me. In this project I have implemented some things that could be used in future projects, for example the **Router** or the **Input Validation** system. Only 3 modules are used, **bcrypt**, **sqlite driver**, and **go-jwt**, this allows me to create a simple project with only the really necessary modules.

I have implemented just 4 routes, since the aim was to create something simple. These routes are:
- `POST http://localhost:3000/register` Register a new user. It checks if the **username** is already taken.
- `POST http://localhost:3000/login` Returns a JWT if the credentials exist.
- `POST http://localhost:3000/isAuthenticated` Returns a boolean if you are authenticated.
- `POST http://localhost:3000/protectedRoute` Returns a message if you are authenticated.

# Running
It's quite easy! You just need to pass one env variable `SECRET_JWT_KEY`. So you can just run this:
```bash
SECRET_JWT_KEY=secret go run ./app/
```

