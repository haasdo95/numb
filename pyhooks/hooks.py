import torch
import os
from contextlib import redirect_stdout, redirect_stderr
import torch.onnx
import json

GRAPH_FD = 3
PARAM_FD = 4
STATE_DICT_FD = 5

def numb_graph(model, dummy_input):
    mode = os.getenv("NUMB_MODE")
    if mode == None:  # a nop!
        print("NO OP!")
        return
    writer_pipe = os.fdopen(GRAPH_FD, 'w') # write-end of the pipe
    with redirect_stdout(writer_pipe):
        torch.onnx.export(model, dummy_input, ".nmb/.tmp", verbose=True)

def numb_param(params):
    mode = os.getenv("NUMB_MODE")
    if mode != "TRAIN":  # nop
        print("NO OP!")
        return
    writer_pipe = os.fdopen(PARAM_FD, 'w') # write-end of the pipe
    writer_pipe.write(json.dumps(params))

def numb_state_dict(state_dict):
    mode = os.getenv("NUMB_MODE")
    if mode != "TRAIN":  # nop
        print("NO OP!")
        return
    writer_pipe = os.fdopen(STATE_DICT_FD, 'wb') # write-end of the pipe
    torch.save(state_dict, writer_pipe)
