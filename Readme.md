# About

Vikingo Satellite is a cross-platform utility for exploiting vulnerabilities without explicitly displaying the result like:

- Out-of-band remote code execution
- Out-of-band SQL Injection
- XXE
- Server-side request forgery
- etc

Key features:

- one file ~20mb (no need to link 3th party and dependencies)
- user-friendly UI
- HTTP module supports automatic renewal of Letsencrypt TLS certificates
- cross-platform
- integration with Vikingo Engine

Built-in modules:

- dns
- http
- ftp
- tcp

# How to start

You can download builded release at here: https://github.com/vikingo-project/vsat/releases or use docker image.

For example:

```
docker run --rm -ti -v -p 1025:1025 vkngo/satellite ./vsat64

```

Or use persistent DB storage

```
docker run --rm -ti -v /your/persistent/directory:/app/storage -p 1025:1025 -p 53:53/udp -p 80:80 -p 443:443 -p 21:21 -p 60000:60000 -p 60001:60001 vkngo/satellite ./vsat64 -db /app/storage/dbname.db
```

where `/your/persistent/directory` is a directory on the host system.

# Screenshots and demo

<details>
  <summary>Interactions page</summary>
  <img src="https://static.vikingo.org/images/satellite-1.jpg"/ >
</details>
<details>
  <summary>Share folder via HTTP</summary>
  <a href="https://www.youtube.com/watch?v=uPK0ltzT5o4"><img src="https://i9.ytimg.com/vi/uPK0ltzT5o4/mq2.jpg?sqp=CISP6YoG&rs=AOn4CLDUE2vpBcxG2C-mGOxpIRcoj2G2bg" /></a>
</details>
