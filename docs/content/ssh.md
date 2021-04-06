---
title: "SSH"
date: 2020-01-01T00:00:00-00:00
draft: false
---
# SSH Server

The SSH server requires that the host has SSH keys, so generate SSH keys if not already:

```bash
$ ssh-keygen
```

Check keys were generated:

```bash
$ ls ~/.ssh
id_rsa id_rsa.pub
```

Run SSH server:

```bash
cointop server -p 2222
```

If the host SSH keys live elsewhere, specify the location:

```bash
cointop server -p 2222 -k ~/.ssh/some-dir/id_rsa
```

SSH into server to see cointop:

```bash
ssh localhost -p 2222
```

The cointop SSH server will use the client's public SSH key as the identifier for persistent config by default. You may change it to use the username instead:

```bash
cointop server -p 2222 --user-config-type=username
```

SSH'ing into server with same username will use the same respective config now:

```bash
ssh alice@localhost -p 2222
```

Pass arguments to cointop on SSH server by using SSH `-t` flag followed by cointop command and arguments. For example:

```bash
ssh localhost -p 2222 -t cointop --colorscheme synthwave
```

## Using docker to run SSH server:

```bash
docker run -p 2222:22 -v ~/.ssh:/keys --entrypoint cointop -it cointop/cointop server -k /keys/id_rsa
```

cointop server writes the client config to `/tmp/cointop_config` within the container, so to make it persistent in host attach a volume. The following example will to write the cached config to `~/.cache/cointop` on the host:

```bash
docker run -p 2222:22 -v ~/.ssh:/keys -v ~/.cache/cointop:/tmp/cointop_config --entrypoint cointop -it cointop/cointop server -k /keys/id_rsa
```

## SSH demo

```bash
ssh cointop.sh
```
