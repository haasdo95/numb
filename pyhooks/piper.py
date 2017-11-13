import os

f = os.fdopen(3, "w")

print("hello!")

f.write("hello from py!")
f.close()