import sys
sys.path.append('/Users/dwd/go/src/github.com/user/numb')

import torch
from torch import nn
from torch.nn import functional as F
from torch.autograd import Variable
from pyhooks.hooks import numb_model

dummy_input = Variable(torch.ones((1, 3, 32, 32)))

class LeNet(nn.Module):
    @numb_model(dummy_input)
    def __init__(self):
        super(LeNet, self).__init__()
        self.conv1 = nn.Conv2d(3, 6, 5)
        self.fc1   = nn.Linear(1176, 120)
        self.fc3   = nn.Linear(120, 10)

    def forward(self, x):
        out = F.relu(self.conv1(x))
        out = F.max_pool2d(out, 2)
        out = out.view(out.size(0), -1)
        out = F.relu(self.fc1(out))
        out = self.fc3(out)
        return out

if __name__ == '__main__':
    net = LeNet()
    print(net)
    