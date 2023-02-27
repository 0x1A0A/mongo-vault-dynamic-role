#!/bin/bash

vault auth enable -path="mongouser" userpass

vault write auth/mongouser/users/mongo password="v-pwd" policies=my-mongo