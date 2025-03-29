#!/usr/bin/python

import socket
import os, os.path

if os.path.exists("/tmp/unix.sock"):
  os.remove("/tmp/unix.sock")

server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
server.bind("/tmp/unix.sock")
while True:
  server.listen(1)
  conn, addr = server.accept()
  conn.setblocking(True)
  print(conn.recvmsg(2048))
