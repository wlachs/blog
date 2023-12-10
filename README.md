# wlchs/blog

![CircleCI](https://dl.circleci.com/status-badge/img/circleci/TGBigTiT6XaXL1AJ6eQozq/9bfPtNbp6S3XPXo2a15wNy/tree/main.svg?style=shield&circle-token=d045b5976dddbd61fa4267e01ba2abb4feca1659)

A lightweight blog engine implemented in Go.

---

# Setup

Follow the guide to deploy your very own blog engine!

## Prerequisites

For a successful deployment, you need Docker.
You can use [Docker Desktop](https://docs.docker.com/desktop/) or [Colima](https://github.com/abiosoft/colima) if you prefer a
command-line-only solution.

## Configuration

To customize the blog engine, you must edit the configuration files.
These contain, among other essential settings, the primary user's name.
The configuration files are at [build/env](./build/env).
**Do NOT upload these files to your version control.**

I highly recommend you change the **highlighted** properties.

**core.env:**

| Key                  | Default | Description                                                         |
|----------------------|---------|---------------------------------------------------------------------|
| **JWT_SIGNING_KEY**  | -       | This should be a strong password for signing authentication tokens. |
| **DEFAULT_USER**     | -       | Name of the primary user. Change this to your name.                 |
| **DEFAULT_PASSWORD** | -       | Primary user's password.                                            |
| GIN_MODE             | RELEASE | Leave in on "RELEASE" unless you know what you're doing.            |

**shared.env:**

| Key            | Default    | Description                                                                                         |
|----------------|------------|-----------------------------------------------------------------------------------------------------|
| MYSQL_USER     | blog_admin | Database username. There is no need to change if you use the preconfigured MySQL docker container.  |
| MYSQL_PASSWORD | password   | Database password. There is no need to change if you use the preconfigured MySQL docker container.  |
| MYSQL_DATABASE | blog       | Database schema. There is no need to change if you use the preconfigured MySQL docker container.    |
| MYSQL_HOST     | db         | Database hostname. Change this if you use your database instead of the one in the docker container. |
| MYSQL_PORT     | 3306       | Database port. Change this if you use your database instead of the one in the docker container.     |

**db.env:**

| Key                 | Default | Description                                                                      |
|---------------------|---------|----------------------------------------------------------------------------------|
| MYSQL_ROOT_PASSWORD | -       | Database root password. There is no need to provide it if you use your database. |

## Deployment

After successfully customizing your configuration files, there is only one more step: deployment.

```sh
docker compose build
docker compose up
```

## For contribution and development

If you'd like to run the blog engine in developer mode to test it or contribute, there are a few differences.

First, you need a database. You can deploy a MySQL database in a Docker container like in a "real" release.
Just make sure you remember the username and the password.

Then, you must set your machine's environment variables as in the configs above.
I recommend using a tool such as [direnv](https://direnv.net).
In this case, you can set the variables the following way:

```sh
export JWT_SIGNING_KEY=SuperSecret
export MYSQL_ROOT_PASSWORD=root
export MYSQL_USER=blog_admin
export MYSQL_PASSWORD=password
export MYSQL_DATABASE=blog
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export PORT=8080
export DEFAULT_USER=TestUser
export DEFAULT_PASSWORD=Test1234
```

Of course, you should change the values to match the ones you have used for the database deployment.

Last but not least, you have to install Go.
If you are using macOS and have Homebrew installed, you can install Go by running the following command:

```sh
brew install go
```

Now that everything is ready, you can navigate to [cmd/blog](./cmd/blog) and run the following command:

```sh
go run .
```

Happy coding!

# Testing

To ensure the stability of the blog engine and that new features don't accidentally break existing ones, I've decided to implement unit
tests. You can follow the current state of test coverage on various software components in the table below.

| Component        | Coverage (%) | State              |
|------------------|--------------|--------------------|
| **Controllers**  |              |                    |
| AuthController   | 100%         | :white_check_mark: |
| PostController   | 100%         | :white_check_mark: |
| UserController   | 100%         | :white_check_mark: |
| **Services**     |              |                    |
| PostService      | 100%         | :white_check_mark: |
| UserService      | 100%         | :white_check_mark: |
| **Repositories** |              |                    |
| PostRepository   | 0%           | :x:                |
| UserRepository   | 0%           | :x:                |
| **Utils**        |              |                    |
| AuthUtils        | 100%         | :white_check_mark: |
| TokenUtils       | 100%         | :white_check_mark: |
