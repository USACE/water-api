"""
Utility to keep the container running if one needs to run
processor.py manually for testing/debugging.
"""
import time

while True:
    print("Init script sleeping...")
    time.sleep(30)
