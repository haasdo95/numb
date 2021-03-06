# TODO #1: Refactor every occurrence of "test" into "trial", since the term "test" has a commonly accepted meaning in machine learning and doesn't perfectly match what `numb test` really does(tuning or cross-validation).


# `Numb` User Manual

# 1. Command Line Tools

## `numb init`
`numb init` initiaizes the project for train-trail cycle management.

After running the command under your project root, you should see:

- `nmb.json` 
	- A json file to configure which command will be run for train/trial, together with a few other options to fill out. 

- a new git branch named "numb"
	- This will become the home to all dirty commits generated by `numb`. I know some people could be really weird about the commits in their projects.

- a hidden directory named `.nmb`
	- Well since it is hidden you can not technically see it. 
	- For the curious, this directory is where `numb` stores stuff, which is primarily the `state dictionaries` generated by training.

## `numb deinit`

***Caveat: this command is primarily prepared for the ease of testing and is subject to removal in near future.***

Basically everything created by running `numb init` will be destroyed by running `numb deinit`, like the destined whirlwind of Macondo.

## `numb train`
Ya I know you've been waiting for this.

Basically whatever command you put in the "train" section of `nmb.json` will be run in train mode.

`numb` does the following when `numb train` is run ***(important)***:

- The computation graph of your architecture will be retrieved (thanks to the magic of `onnx`).
- The hyperparameters you registered with `numb_param` will be retrieved. If this architecture has been trained with the same set of hyperparameters before (which means you're redoing your work!), you'll receive a prompt.
- By the time your training script finishes running, the `state dictionary` will be extracted and automatically saved. So, yeah, no more boilerplate.
- A commit will be made on `numb` branch, so that you could come back and admire your great work any time you want.

## `numb trial`
Basically whatever command you put in the “trial” section of nmb.json will be run in trial mode.

`numb` does the following when `numb trial` is run ***(also important)***:

- The computation graph of your architecture will be retrieved (thanks to the magic of `onnx`). If *no* model has ever been trained with the architecture, your trial script will be shutdown.
- You will be prompted to choose a specific model trained with a specific combination of hyperparameters.
- The trial result (accuracy, etc.) you registered with `numb_trial_result` will be retrieved to fill out the untrialed entry.

## `numb trial -all`

By running `numb trial -all`, you could trial all the untrialed models associated with current architecture.

## `numb queue init <filename>`

***caveat: the whole `numb queue` business may be integrated into `numb train` as an extra cmdline flag like `numb train -queue`***

Create a queuefile named as `<filename>`.

The queuefile is basically a json file with a list of hyperparameter combinations.

## `numb queue run <filename>`

By running this command,

- Your train script will be run for several times.
- For each time, a different hyperparameter combination will be "injected" into your train script, ***as long as you remember to put `numb_queue(globals())` in your training script before your hyperparameters are referenced***

## `numb list`

`numb list` lists the models, together with helpful information.

Whenever you are lost in the `numb` workflow, run `numb list` to get back on track.

Listed models will be grouped by architecture, with information of hyperparameters associated with the architecture listed below.

## `numb revert <ID>`

Well, for now the `<ID>` is nothing but a timestamp.

This command allows you to jump back to a previous stage of the project associated with a specific training record.

This is the command to run if you want to figure out why a specific architecture or a hyperparameter combination works so well.

## `numb report <ID>`

This command takes out the `state dictionary` associated with `<ID>` from the messy `.nmb` folder and put it in a new directory named as "report-`<ID>`"

# 2. Python Hooks

## `numb_model`

A decorator to the neural network constructor.

Basically, write `@numb_model(<dummy_input>)` above your `__init__`

- What is the dummy input?
	- A dummy input is a `Variable` with the dimension of your input. You could create it with `torch.ones()` or `torch.zeros()`

- Why do we need it?
	- I know this is a pain. But unlike TensorFlow, PyTorch adopts dynamic computation graph, which allows for easier debugging and prototyping. As a result, we need a full run through your neural net architecture to retrieve the computation graph, necessitating a dummy input.

## `numb_param`

Register hyperparameter combination with it.

`numb_param` is a function that takes in a dictionary, which specifies the hyperparameter values.

## `numb_queue`

A function used to "inject" hyperparameters into your training script. 

Put `numb_queue(globals())` before your hyperparams are referenced.

- TODO: Why do we need to put `globals()` here?

## `numb_trial_result`

A function to retrieve your trial result by taking in a dictionary. 