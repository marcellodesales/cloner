# Make local binary

```
make local
```

# Private Repos - git@host:org/repo

1. Create a test SSH key
  - Convert the OPENSSH to RSA
2. Add the public key as a `Deploy Key` 
  - `Settings -> Deploy Keys`
  - Make sure to create a deploy key *with clone permission only*!
3. Run the container providing the private key as param
4. References

## Create a test SSH key

* Just create a test key with empty passphrase

```
$ ssh-keygen -t rsa -b 4096 -C "test@cloner.github.com" -N ""
Generating public/private rsa key pair.
Enter file in which to save the key (/Users/marcellodesales/.ssh/id_rsa): id_cloner_test
Your identification has been saved in id_cloner_test.
Your public key has been saved in id_cloner_test.pub.
The key fingerprint is:
SHA256:XP89ufanVLXe63vYJSrRYGzWX0cAV+1TJAhgffs3lUw test@cloner.github.com
The key's randomart image is:
+---[RSA 4096]----+
|        oo...o+++|
|       .  . o. E+|
|         ..o .oo+|
|       . .*.o  +B|
|        S+ o.o o=|
|          . ..=+=|
|           . .oO=|
|          . ....B|
|           .  +**|
+----[SHA256]-----+
```

* 2 files are created: `id_cloner_test` and `id_cloner_test.pub`.

```
$ cat id_cloner_test
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAACFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAgEAt/DB1OZFiWqwHXjRItqFnxOIZ2+st6OTkF+b7OGD8vnMdm058s1K
...
...
PL1esI1hYGI9Z6II+g3iQwxkszTaBRs6vLPBsKqkcSH4CvlvOPbSDJwPgVvvcgRTZPF7rZ
u31hhB/MXHnN4yDFjiTqFHZWQgTlUWTvrqdZKkO8XGeisu0Z9vzn7XHjDdgN69jrjPAANh
UTASRHJrLmqzAAAAFnRlc3RAY2xvbmVyLmdpdGh1Yi5jb20BAgME
-----END OPENSSH PRIVATE KEY-----

$ cat id_cloner_test.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC38MHU5kWJarAdeNEi2oWfE4hnb6y3o5OQX5vs4YPy+cx2bTnyzUpJZZ1mQmRupTIz+tssp12Co/HFxAOd9KNqIpdN42FuFVAUO1waGcIqRELfF+/kYMJPk4ZxKMWe+gJ9KHqPGlFeUlsmVoGERYuHFUOIXptKLqxrkTh+fEVzT7eUyY8DXnjOlggjoBCm9VTcHvlETndlJu6LH9CtXMuP18noIXAB9mqWIe4dqf1jdXbOANiCLYId8H84mnjjML6aFgE4JO3D/WBQFwI5Vr1WNL0g45NbJ+j3MKVm2MXQAzd4RE2mH6Rv5wr4cLzvHTfy2x+Cm6fq5kml3N4uku+4BBVWXVCxd+AXf42z5/kmOPQwo58Kk++CQhnaP49UroXW8NmQ6zMUsB5mZqypQ4NPLrLCXQcK4jbUSjG3fLquEv3pmzXSiqQ6618SwRKDsC0qIGanDsZmKthm3PabG2/qEzkDES6axFrJAQQGSGmScJP+ZSxNyvsiSzbQqRfPMcHYOW9OmGpGvGK4SO+46alcBhEUhqGyZBvKJdWfCG8qeEZ0M9uyLmZrMJ6i+aMpNReFmhAx43sGd6SpGHhb4uGUkGiSABLD51ZVW/I1rOIEytPImBVZa+5wkwK3JmFs1savVUqnX7Q0SH9N9JL7KSih1W/jCYHvOw8h0bro4JmUOQ== test@cloner.github.com
```

* Convert the key from `OPENSSH` to RSA

```
$ ssh-keygen -p -m PEM -f id_cloner_test
Key has comment 'test@cloner.github.com'
Enter new passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved with the new passphrase.
```

* The file is updated in place

```
$ cat id_cloner_test
-----BEGIN RSA PRIVATE KEY-----
MIIJKAIBAAKCAgEAt/DB1OZFiWqwHXjRItqFnxOIZ2+st6OTkF+b7OGD8vnMdm05
8s1KSWWdZkJkbqUyM/rbLKddgqPxxcQDnfSjaiKXTeNhbhVQFDtcGhnCKkRC3xfv
5GDCT5OGcSjFnvoCfSh6jxpRXlJbJlaBhEWLhxVDiF6bSi6sa5E4fnxFc0+3lMmP
...
...
GmiYxMLb1gQYVH+R9wxIb0LaBl+XLfjKjU8Z+aW0lVVXcAIgCIwG0qdL5CiY0rG4
coclH2OmlSN2G+fS8jK/dcAHu7oQZptFUZ9Q98PuVsi9csTfxmN/y1mcq5o=
-----END RSA PRIVATE KEY-----
```

## Add a Deploy Key

* https://developer.github.com/v3/guides/managing-deploy-keys/#deploy-keys
  * https://github.com/marcellodesales/cloner/settings/keys/new
  
Use the public key so that CI systems can clone it

> **ATTENTION**: Make sure the key is listed as `Read-only`.

## Run the container to clone

* Execute the command

```
$ docker run -ti -v $(pwd):/data marcellodesales/cloner:20.09.7 git \
         --repo git@github.com:marcellodesales/unmazedboot.git -v debug -k /data/id_cloner_test
DEBU[2020-09-18T02:17:38Z] config.git.cloneBaseDir=/root/cloner
INFO[2020-09-18T02:17:38Z] Cloning into '/root/cloner/github.com/marcellodesales/unmazedboot'
DEBU[2020-09-18T02:17:38Z] Attempting to clone repo 'git@github.com:marcellodesales/unmazedboot.git' => '/root/cloner/github.com/marcellodesales/unmazedboot'
DEBU[2020-09-18T02:17:38Z] Authenticating using the key
Enumerating objects: 16, done.
Counting objects: 100% (16/16), done.
Compressing objects: 100% (16/16), done.
Total 144 (delta 9), reused 2 (delta 0), pack-reused 128
INFO[2020-09-18T02:17:41Z] Done...
INFO[2020-09-18T02:17:41Z]
/root/cloner/github.com/marcellodesales/unmazedboot
└── .env
└── CHANGELOG
└── LICENSE
└── README.md
└── builder
│   ├── gradle.Dockerfile
│   ├── maven.Dockerfile
```

## References

* https://stackoverflow.com/questions/44269142/golang-ssh-getting-must-specify-hoskeycallback-error-despite-setting-it-to-n/63308243#63308243
* https://skarlso.github.io/2019/02/17/go-ssh-with-host-key-verification/
* https://github.com/src-d/go-git/issues/637#issuecomment-543015125

> NOTE: FOR VERIFICATION with agent-based execution

Needs to execute the agent https://github.com/src-d/go-git/issues/550#issuecomment-323075887

In order to test, run a docker container as follows:

1. Start an ssh-agent

In order to avoid errors like the following:

```
ERRO[2020-09-17T21:39:33Z] can't clone the repo at 'github.com/marcellodesales/unmazedboot': 
       error creating SSH agent: "SSH agent requested but SSH_AUTH_SOCK not-specified"
```

* Test in Container

```
$ docker run -ti -v $(pwd):/data --entrypoint sh alpine/git
/git # cd /data/
/data # eval `ssh-agent`
Agent pid 19
/data # ssh-add ~/.ssh/id_cloner_test
Identity added: /root/.ssh/id_cloner_test (/root/.ssh/id_cloner_test)
```

2. If the known_hosts is empty, then errors will show

```
ERRO[2020-09-17T22:26:53Z] can't clone the repo at 'github.com/marcellodesales/unmazedboot': 
       ssh: handshake failed: knownhosts: key is unknown
```

3. Copy the known hosts file to verify

```
/data # cp known_hosts ~/.ssh/known_hosts
```

For this reason, validation on the ssh key is disabled when the file does not exist

4. If the file exists, there's no verification.

Build local and test.