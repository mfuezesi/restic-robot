#!/usr/bin/env python3

import random
import os

random.seed(42)

def recurse(base_path: str, more_depth: int):
    for i in range(100):
        fn = "d{}".format(i)
        path = os.path.join(base_path, fn)
        if more_depth > 1:
            os.mkdir(path)
            recurse(path, more_depth - 1)
        else:
            with open(path, "wb") as f:
                length = max(1, int(random.expovariate(1.0/4000)))
                data = random.getrandbits(length*8).to_bytes(length, "little")
                f.write(data)

os.mkdir("test")
recurse("test", 3)