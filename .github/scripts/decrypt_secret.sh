#!/bin/sh

# Decrypt the file
# --batch to prevent interactive command
# --yes to assume "yes" for questions
gpg --quiet --batch --yes --decrypt --passphrase="${ID_CLONER_TEST_PASSPHRASE}" --output $2 $1
