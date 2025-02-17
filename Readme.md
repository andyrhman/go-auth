# Go Auth

<h1 align="center">
  <a href="https://gofiber.io">
    <picture>
      <source height="125" media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo-dark.svg">
      <img height="125" alt="Fiber" src="https://raw.githubusercontent.com/gofiber/docs/master/static/img/logo.svg">
    </picture>
  </a>
</h1>
<p align="center">
  <em><b>Fiber</b> is an <a href="https://github.com/expressjs/express">Express</a> inspired <b>web framework</b> built on top of <a href="https://github.com/valyala/fasthttp">Fasthttp</a>, the <b>fastest</b> HTTP engine for <a href="https://go.dev/doc/">Go</a>. Designed to <b>ease</b> things up for <b>fast</b> development with <a href="https://docs.gofiber.io/#zero-allocation"><b>zero memory allocation</b></a> and <b>performance</b> in mind.</em>
</p>

This is my first time doing a golang project using go fiber, back before when i started to do this i was thinking which API framework has the fastest response and lowest latency on google and after find some website showing the benchmark i found out go fiber rank 3rd the fastest API framework, this leads me to try studying go fiber in 2025 after several months of my busy independent study on Bangkit.

## Installation

```bash
go mod init nameof/yourproject
go get -u github.com/gofiber/fiber/v3
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/google/uuid
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto
go get -u gorm.io/driver/postgres
go get -u gorm.io/gorm
```
