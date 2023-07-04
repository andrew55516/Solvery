# Solvery


## Features
___
- Create user
- Update user credit
- Get user by email
- List users
- List user entries
- List all entries
- Solve task for credit
- Send emails with task answer and credit balance

## Tech
___
- Built in [go](https://go.dev/) version 1.20
- Uses [docker](https://www.docker.com/) 
- Uses [Postgres](https://www.postgresql.org/)
- Uses [gin-gonic](https://github.com/gin-gonic/gin)
- Uses [sqlc](https://github.com/kyleconroy/sqlc)
- Uses [email](https://github.com/jordan-wright/email)
- Uses [testify](https://github.com/stretchr/testify)
- Uses [viper](https://github.com/spf13/viper)

## API
___
You can import postman [collection](https://github.com/andrew55516/Solvery/tree/postman/postman/collections) to test API

## Start
___
Download the repository and run following command to run app with docker-compose:
```
$ make dockerserver
```

## Details
___
- Must be some identification for each user: let it be the user email
- Need I store the entries for each operation with user credit? I guess, it has a point
- Should it be debt or credit for users balance? Users can refill balance, not just pay off a debt. So, let it be a credit with minimal value -1000, and it also can be positive
- How much costs task solution? In my mind, it should depend on input size and algorithm complexity. For this reason, solution for task1 costs n, where n is len(input), such as time and space complexity of my solution is _O(n)_ 
- I implemented [workflow](https://github.com/andrew55516/Solvery/tree/postman/.github/workflows) for build and testing
- Code coverage with tests is above 77%