#!/bin/bash

set -x

# an invalid password should fail
if [[ ${STIG} == true ]]; then
    echo "tooshort12@" | sudo pwscore && exit 1
else
    echo "tooshort1@" | sudo pwscore && exit 1
fi
echo "password123456@J" | sudo pwscore && exit 1
echo "passSDWword@@@@J" | sudo pwscore && exit 1
echo "passSDWword1111J" | sudo pwscore && exit 1
echo "lowerrrr12case@" | sudo pwscore && exit 1
echo "UPPERRR12CASE@" | sudo pwscore && exit 1

# a valid password should succeed
echo "passSDWword1232rdw#@" | sudo pwscore || exit 1

# validate password age settings
if [[ ${STIG} == true ]]; then
    grep -E '^PASS_MAX_DAYS 60$' /etc/login.defs || exit 1
    grep -E '^PASS_MIN_DAYS 1$' /etc/login.defs || exit 1
    grep -E '^UMASK 077$' /etc/login.defs || exit 1
else
    grep 'PASS_MAX_DAYS 90' /etc/login.defs || exit 1
    grep 'PASS_MIN_DAYS 7' /etc/login.defs || exit 1
fi
grep 'INACTIVE=30' /etc/default/useradd || exit 1
