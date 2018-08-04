#!/bin/sh

. env
. env-test

export $(cut -d= -f1 env)
export $(cut -d= -f1 env-test)

./triagebot
