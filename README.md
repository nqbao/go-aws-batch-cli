# AWS Batch CLI

AWS Batch CLI is a command line utility to interact with [AWS Batch](https://aws.amazon.com/batch/). The motivation is to simplify the workflow of developing and debugging a job, which is currently super painful due to the UI.

## How to use

Run `make build` to build the command line. Use `help` subcommand to explore all the options.

**Run a job**:

Use `run` subcommand to run a batch job

```
./build/aws-batch-cli run -j job-definition -q queue-name -p param1=value1 -p param2=value
```

You can also use `-f` to follow the log stream

```
./build/aws-batch-cli run -j job-definition -q queue-name -p param1=value1 -p param2=value -f
```

To override the container command

```
./build/aws-batch-cli run -j job-definition -q queue-name -p param1=value1 -p param2=value -f -- my container command
```

**View log of an existing job**

```
./build/aws-batch-cli log -i abcdef
```

## TODO

 - [ ] Describe a job
 - [ ] Cancel a job
 - [ ] Rerun a job
 - [ ] Automate job building and packaging
