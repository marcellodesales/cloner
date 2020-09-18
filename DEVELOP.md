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
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
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
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
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