import sys
import os

sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
# print(sys.path)

from job import JobService
from thrift import Thrift
from thrift.transport import TSocket, TTransport
from thrift.protocol import TBinaryProtocol, TCompactProtocol


try:
    t = TSocket.TSocket('localhost', 9000)
    transport = TTransport.TBufferedTransport(t)
    protocol = TCompactProtocol.TCompactProtocol(transport)
    client = JobService.Client(protocol)
    transport.open()
    res = client.ListJob()
    print(res.Job[0].Name)
    print(res.Job[0].Command)
    print(res.Job[0].Expr)
    transport.close()
except Thrift.TException as e:
    print(e.message)
