import sys
sys.path.append('/Users/dwd/go/src/github.com/user/numb')

from torch import nn
from torch.nn import functional as F

class LeNet(nn.Module):
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
    