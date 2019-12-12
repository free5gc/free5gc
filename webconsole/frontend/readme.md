# free5GC Web Console

## Environment Setup

Install yarn
```bash
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt update && sudo apt install yarn
```

Install required packages
```bash
yarn install
```

## Run the Web Server
```bash
yarn start
```
