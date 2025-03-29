#!/usr/bin/python

import socket

server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
server.connect("/tmp/unix.sock")
server.send(bytes(str(input("Whatcha wanna say? ")), encoding='utf8'))
