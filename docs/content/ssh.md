---
title: "SSH"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# SSH Server

Run SSH server:

```bash
cointop server -p 2222
```

SSH into server to see cointop:

```bash
ssh localhost -p 2222
```

SSH demo:

```bash
ssh cointop.sh
```

Passing arguments to SSH server:

```bash
ssh cointop.sh -t cointop --colorscheme synthwave
```

Using docker to run SSH server:

```bash
docker run -p 2222:22 -v ~/.ssh:/keys --entrypoint cointop -it cointop/cointop server -k /keys/id_rsa
```
