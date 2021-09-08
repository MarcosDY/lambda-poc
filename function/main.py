#!/usr/bin/env python3

import os
import base64

def pop_handler(event, context):
    f = open("/tmp/svid.pem", "r")
    cert = f.read()
    print(cert)

    f = open("/tmp/bundle.pem", "r")
    bundle = f.read()

    return { 
        'cert' : cert,
        'bundle': bundle
    }
if __name__ == "__main__":
    pop_handler(None, None)
