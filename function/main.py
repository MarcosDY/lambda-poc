#!/usr/bin/env python3

import os
import base64

def pop_handler(event, context):
    f = open("/tmp/svid.0.pem", "r")
    s = f.read()
    print(s) 

if __name__ == "__main__":
    pop_handler(None, None)

