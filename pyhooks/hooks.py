import torch
from torch import nn
import os
from contextlib import redirect_stdout, redirect_stderr
import torch.onnx
import json
import signal
import atexit
import sys

GRAPH_FD = 3
PARAM_FD = 4
STATE_DICT_FD = 5
INTERACT_FD = 6

def numb_graph(model, dummy_input):
    mode = os.getenv("NUMB_MODE")
    if mode == None:  # a nop!
        print("NO OP FOR GRAPH!")
        return
    writer_pipe = os.fdopen(GRAPH_FD, 'w') # write-end of the pipe
    with redirect_stdout(writer_pipe):
        torch.onnx.export(model, dummy_input, ".nmb/.tmp", verbose=True)
    writer_pipe.close()

def numb_param(params):
    mode = os.getenv("NUMB_MODE")
    if mode != "TRAIN":  # nop
        print("NO OP FOR PARAM!")
        return
    writer_pipe = os.fdopen(PARAM_FD, 'w') # write-end of the pipe
    writer_pipe.write(json.dumps(params))

def numb_state_dict(model: nn.Module):
    mode = os.getenv("NUMB_MODE")
    if mode != "TRAIN":  # nop
        print("NO OP FOR STATE DICT!")
        return
    writer_pipe = os.fdopen(STATE_DICT_FD, 'wb') # write-end of the pipe
    torch.save(model.state_dict(), writer_pipe)

def numb_test(model: nn.Module):
    """
    this one will block and wait for StateDictFileName
    :return:
    """
    mode = os.getenv("NUMB_MODE")
    if mode != "TEST":
        print("NO OP FOR TEST!")
        return
    def handle_usr2(): # exit on sigusr2
        print("SHUTTING DOWN")
        sys.exit(1);
    signal.signal(signal.SIGUSR2, handle_usr2)
    reader_pipe = os.fdopen(INTERACT_FD, "r") # wait for user choice of state dict
    os.kill(os.getppid(), signal.SIGUSR1)
    state_dict_filename = reader_pipe.read()
    print("State Dict Filename: ", state_dict_filename)
    with open(state_dict_filename, "rb") as sdf: # read in state dict and load
        model.load_state_dict(sdf)


def numb_model(dummy_input):
    def actual_decorator(initfunc):
        def wrapped(*args, **kwargs):
            initfunc(*args, **kwargs)
            numb_graph(args[0], dummy_input)
            numb_test(args[0])
            atexit.register(numb_state_dict, [args[0]])
        return wrapped
    return actual_decorator