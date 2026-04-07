#!/usr/bin/env bash

run() {
    echo -e "\e[33m\$ $*\e[0m"
    read -r
    "$@"
    echo ""
}

clear
echo "imscli Demo"
echo ""
read -r

# 1. DCR
run ./imscli dcr register \
    --url https://ims-na1-stg1.adobelogin.com \
    --clientName "My App" \
    --redirectURIs "http://localhost:8888" \
    --scopes "openid AdobeID"

read -r -p "client_id: " CLIENT_ID
read -r -p "client_secret: " CLIENT_SECRET

# 2. Authz
echo -e "\e[33m\$ ./imscli authz user --url https://ims-na1-stg1.adobelogin.com --clientID $CLIENT_ID --clientSecret $CLIENT_SECRET --scopes \"openid AdobeID\" --organization \"96D656605F9846070A494236@AdobeOrg\"\e[0m"
read -r
ACCESS_TOKEN=$(./imscli authz user \
    --url https://ims-na1-stg1.adobelogin.com \
    --clientID "$CLIENT_ID" \
    --clientSecret "$CLIENT_SECRET" \
    --scopes "openid AdobeID" \
    --organization "96D656605F9846070A494236@AdobeOrg")
echo "$ACCESS_TOKEN"
echo ""

# 3. OBO
run ./imscli obo \
    --url https://ims-na1-stg1.adobelogin.com \
    --clientID "imscli" \
    --clientSecret "s8e-QZ9CnNQq2W6J_suZdqNd__BaaXokfZVq" \
    --accessToken "$ACCESS_TOKEN" \
    --scopes "openid AdobeID"

echo -e "\e[32mDone!\e[0m"
