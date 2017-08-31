Nextpass generates cryptographically secure passwords.
It can either be used on the command line, or imported by other go projects.
Whenever you need to create your next password, use nextpass.

```sh
go get -u github.com/n0ot/nextpass
# change to github.com/n0ot/nextpass
make install
```

## Usage
```sh
nextpass [options]
```

    -A, --additional    read additional characters from standard input, encoded in UTF-8; newline characters will NOT be ignored
    -D, --digits        include digits 0-9
    -l, --length uint   length of resulting password (default 64)
    -L, --lower         include lowercase letters a-z
    -n, --no-newline    Don't print a newline after the password
    -r, --random-source string   specify a file to be used as an alternate source of randomness. Don't use this unless you know what you're doing.
    -S, --special       include special characters, which are the printable ascii characters excluding letters, digits, and the space
    -U, --upper         include uppercase letters A-Z
    -v, --verbose       print more information, in addition to the generated password

If the included characters are not enough,
pass your favorite foreign characters or emojis into standard input.

Examples:

```sh
nextpass -l 32 -LUDS
```
generates a password with length 32, including
lowercase and uppercase letters, digits, and special characters.

```sh
echo -n ABCDEF | nextpass -DA
```
generates a 64 digit hexadecimal string (256 bits).

## Security
[crypto/rand](https://godoc.org/crypto/rand) is used as a source of random entropy by default.
The required number of bytes to generate a password are read all at once,
and the result is encoded using the chosen set of characters.
