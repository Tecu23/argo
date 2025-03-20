# Using NNUE for evaluation

## NNUE

NNUE (Efficiently Updatable Neural Network) is, broadly speaking, a neural network
architecture that takes advantage of having minimal changes in the netwwork inputs
between subsequen evaluations.

NNUE operated on the following principles:

1. The network should have realtively low amount of non-zero inputs.
2. The inputs should change as little as possible between subsequent evaluations.
3. The network should be simple enough to facillitate low-precision inference in
   integer domain.

Following the first principle means that when the network is scaled in size the
inputs must become sparse. Current best architectures have input sparsity in the
order of 0.1%. Small amount of non-zero inputs places a low upper bound on the time
required to evaluate the network in cases where it has to be evaluated in its
entirety. This is the primary reason why NNUE networks can be large while still
being very fast to evaluate.

Following the second principle (provided the first is being followed) creates a
way to efficiently update the network (or at least a costly part of it) instead of
reevaluating it in its entirety. This taks advantage of the fact that a single move
changes the board state only slightly. This is of lower importance than the first
principle and completely optional for the implementations to take advantage of,
but nevertheless gives a measurable improvement in implementations that do care
to utilize this.

Following the third principle allows achieving maximum performance on common
hardware and makes the model especially suited for low-latency CPU inference which
is necessary for conventional chess engines.

Overall the NNUE principles are applicable also to extensive deep networks, but
they shine in fast shallow networks, which are suitable for low-latency CPU
inference without the need for batching and accelerators. The target performance
is million(s) of evaluations per second per thread. This is the extreme use case
that required extreme solutions, and most importantly quantization.
