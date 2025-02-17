import struct
import os

value = 2
filename = "binary_datafile.bin"
with open(filename, "wb") as f:
    while True:
        # Write the 64-bit unsigned integer in little-endian format
        f.write(struct.pack('<Q', value))
        value+=1

        if os.stat(filename).st_size >= 1024 * 1024 * 1024:  # 1GB
            break

print(f"File '{filename}' created. Size: {os.stat(filename).st_size / (1024 * 1024)} MB")
