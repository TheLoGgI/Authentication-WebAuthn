# Authentication Service

## Run dev

**Run Go server**

> `ATLAS_URI=mongodb+srv://username:password@cluster0.37ogv.mongodb.net go run main.go`

**Run React App**

> `cd app && yarn dev`

## Run Production

**Build react app**

This will build the app in the Go servers _static_ directory

> `cd app && yarn build`

Compile the go server to binary

> `go build`

## Generating Public and Private Keys with openssl.exe

source: `http://lunar.lyris.com/help/Content/generating_public_and_private_keys.html`

To perform the following actions for Windows or Linux, you must have OpenSSL installed on your system.

### Generating the Private Key -- Windows

In Windows:

1. Open the Command Prompt (Start > Programs > Accessories > Command Prompt).

2. Navigate to the following folder: C:\Program Files\ListManager\tclweb\bin\certs

3. Type the following: `openssl genrsa -out rsa.private 1024`

4. Press ENTER. The private key is generated and saved in a file named "rsa.private" located in the same folder.

NOTE The number "1024" in the above command indicates the size of the private key. You can choose one of five sizes: 512, 758, 1024, 1536 or 2048 (these numbers represent bits). The larger sizes offer greater security, but this is offset by a penalty in CPU performance. We recommend the best practice size of 1024.

#### Generating the Public Key -- Windows

1. At the command prompt, type the following: `openssl rsa -in rsa.private -out rsa.public -pubout -outform PEM`

2. Press ENTER. The public key is saved in a file named rsa.public located in the same folder.

#### Generating the Private Key -- Linux

1. Open the Terminal.

2. Navigate to the folder with the ListManager directory.

3. Type the following: `openssl genrsa -out rsa.private 1024`

4. Press ENTER. The private key is generated and saved in a file named "rsa.private" located in the same folder.

#### Generating the Public Key -- Linux

1. Open the Terminal.

2. Type the following: `openssl rsa -in rsa.private -out rsa.public -pubout -outform PEM`

3. Press ENTER. The public key is saved in a file named rsa.public located in the same folder.
