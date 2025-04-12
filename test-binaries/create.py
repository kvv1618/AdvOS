import struct
import os
import random

value = random.getrandbits(64)
filename = "binary_datafile_0002.dat"
with open(filename, "wb") as f:
    while True:
        # Write the 64-bit unsigned integer in little-endian format
        f.write(struct.pack('<Q', value))
        value = random.getrandbits(64)

        if os.stat(filename).st_size >= 0.002 * 1024 * 1024 * 1024:
            break

print(f"File '{filename}' created.")
