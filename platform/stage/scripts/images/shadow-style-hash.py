from passlib.hash import sha512_crypt
import getpass

while True:
    password = getpass.getpass("Enter password: ")
    compare = getpass.getpass("Enter password again: ")
    if password != compare:
        print("Passwords are not matching. Try again.")
        continue
    break

hash = sha512_crypt.hash(password)
print("Hash:", hash)
