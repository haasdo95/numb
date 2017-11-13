import sys
sys.path.append('/Users/dwd/go/src/github.com/user/numb')

import torch
from torch import optim
from torch import nn
from torch.autograd import Variable

from pyhooks.hooks import numb_graph, numb_param, numb_state_dict

from demo.lenet import LeNet
model = LeNet()

dummy_inputs = Variable(torch.ones((1000, 3, 32, 32)))
dummy_input = torch.unsqueeze(dummy_inputs[0], dim=0)
dummy_targets = Variable(torch.rand((1000, 10)))

# Register Net!
numb_graph(model, dummy_input)

# set parameters!
learning_rate = 0.001
epoch = 10
batch_size = 20

# Register parameters
numb_param({
    "learning_rate": learning_rate,
    "epoch": epoch,
    "batch_size": batch_size
})

def make_batch(in_data, targets, batch_size):
    for start in range(0, in_data.size()[0], batch_size):
        real_size = min((batch_size, in_data.size()[0] - start))
        yield in_data[start: start + real_size], targets[start: start + real_size]

criterion = nn.MSELoss()
optimizer = optim.SGD(model.parameters(), lr=learning_rate)
for epk in range(epoch):
    batches = make_batch(dummy_inputs, dummy_targets, batch_size)
    for train_exp, target in batches:
        model.zero_grad()
        output = model.forward(train_exp)
        loss = criterion.forward(output, target)
        loss.backward()
        optimizer.step()

# print(model.state_dict())
print("Done Training!")
numb_state_dict(model.state_dict)
