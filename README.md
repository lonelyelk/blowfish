# Blowfish

This is an excercise in implementing [Blowfish](https://www.schneier.com/academic/blowfish/) ([wiki](https://en.wikipedia.org/wiki/Blowfish_(cipher))) algorithm designed by Bruce Schneier.

## Pi

For initializing the P-Array and S-boxes the algorythm uses hexadecimal digits of Pi. It is of course not wise to calculate the digits every time an initializaion needed.
I did however implement the [Bailey–Borwein–Plouffe formula](https://giordano.github.io/blog/2017-11-21-hexadecimal-pi/) for the purpose of the excercise.

## Byte Array Encryption

There is a simple implementation for a byte array encryption (and decryption that assumes the encryption was done this way). First 32 bits (4 bytes) of the first block is the length of the array.
Zero bytes are added in the end to obtain the last whole block.