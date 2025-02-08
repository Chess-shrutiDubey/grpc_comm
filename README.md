# Project: Communication and RPC

The project introduces you to the fundamentals of communication and RPC. The basic idea is to build a library that communicates reliably in the face of failure and gain experience with an RPC package.

## Details

### Language Choice

You can choose to implement the project in either C/C++ or Go. The project is divided into three parts. Whichever language you choose, you must implement all three parts in the same language.  

### Part 1: Measuring & Timing (10% points)

Because the analysis of systems is part of everything you do in systems classes, the first thing you should learn is how to measure how long something takes. On Linux platforms, which we'll be using for this project, you generally use `clock_gettime()` or `gettimeofday()` in C/C++, and in Go you can use `time.Now()` or something similar.

An important aspect of a timer is its precision (or resolution): how small of a time event can be measured with this timer accurately? One way to determine the resolution of the timer is to read the clock value at the start and end of a simple loop. Start with a single loop iteration, then increase the iteration count of the loop until the difference between the before and after samples is greater than zero. Try to get the smallest non-zero positive difference. If a single iteration of a loop takes too much time, try putting simple statements between the two timer calls. Beware compiler optimization, which, if you are not careful, will remove the code in the loop and yield odd results. Record your result.

In this part, you will use your newfound ability to time things to verify [Jeff Dean's numbers](https://gist.github.com/jboner/2841832) you should know. Pick any five entries from that list, measure the time yourself, and make a table to show how your numbers compare against Jeff Dean's numbers.

### Deliverables - Part 1

Comparison table of your measured time and Jeff Dean's time.

### Part 2: Reliable Communications (60% points)

The first real part of the project is to build a reliable communication library on top of raw UDP-based sockets using Go (C/C++). Your communication library should allow two processes to communicate via UDP packets, but it should use a simple timeout-retry mechanism to detect when the receiver has not received a message, and then re-send that message. The sender should keep trying to send the message until it gets an ack from the receiver. Your send code should be blocking, i.e., it should not return until an ack has been received.

After you have this simple layer working, you will measure its performance and reliability characteristics. How much overhead is there to send a message? What is the total round trip time of sending a message and receiving an ack, when running on a single machine, and when running on two separate machines? What is the bandwidth of your library, when sending a large number of max-sized packets between sender and receiver, again when running on the same machine, and when running on different machines? What limits your bandwidth, and how could you do better?

As for reliability, you need to induce controlled message drops to show how your layer works. Arrange for your receiver library code to have an input that tells it what percent of messages it should drop (randomly). If this number is set to 10 percent, for example, your receive-side code should randomly drop 10 percent of messages (naturally); setting the number to 0 makes the layer reliable (no drops). Send a stream of packets, with the reliability percentage set to something non-zero, and measure the round-trip time of each reliable send; in a resulting graph, you should be able to see some high values where time-outs and retries occur. How many different performance regimes result?

When running your experiments, compile your library and test code both with and without optimization enabled (i.e., `-O`) in C/C++; and in Go you may utilize the `-gcflags` flag to enable/disable optimizations. How much difference does this make in your performance results?

### Deliverables - Part 2

1. **Code**:
   - Source code for the reliable communication library using Go (C/C++) UDP packets.
   - Source code for the test cases used to measure performance and reliability.

2. **Performance and Reliability Measurements**:
    A small report detailing:
    - the overhead of sending a message.
    - total round trip time of sending a message and receiving an ack, both on a single machine and on two separate machines.
    - bandwidth measurements when sending a large number of max-sized packets, both on the same machine and on different machines.
    - analysis of what limits the bandwidth and suggestions for improvements.

3. **Reliability Analysis**:
   - Graphs showing the round-trip time of each reliable send with different message drop percentages.

4. **Optimization Analysis**:
   - Comparison of performance results with and without optimization enabled (i.e., using `-gcflags` in Go).

### Part 3: Go gRPC (30% points)

Google has made its RPC system, gRPC, available via open source. We already saw a example of gRPC in the last lab session; and you have probably also built an simple gRPC banking system in the last project. Incase you chose to implement the project in C/C++, you can use the [gRPC C++ library](https://grpc.io/docs/languages/cpp/).

In the last lab, we didn't explicitly utilize Marshal and Unmarshal functions as both the server and client were inside the gRPC framework. gRPC automatically handles the marshalling and unmarshalling of messages. However, if you are working with different RPC frameworks or need to send protobuf messages over non-gRPC channels, you may need to explicitly Marshal and Unmarshal the data. So, it's always an safe practice to utilize Marshal and Unmarshal functions. Guide to Marshal and Unmarshal functions w.r.t. Go protobuf can be found [here](https://pkg.go.dev/google.golang.org/protobuf/proto). Guide to protobuf `proto3` can be found [here](https://protobuf.dev/programming-guides/proto3/).   

What we need now is to measure the overhead of marshalling a message (e.g., packing an item into a protobuf). How long does it take to pack an int, a double, a string (of varying size), or a complex structure on each platform? How confident are you about your measurements?

Now measure round-trip time for a small message, when both client/server are on the same machine, and when on different machines. How long do requests/responses take? Is the first round trip much slower than subsequent ones? How much overhead does using protobuf and RPC take, as compared to your barebones RPC library?

Now measure bandwidth when sending large amounts of data. How quickly can gRPC send large amounts of data from one machine to another? What total bandwidth is achievable, when using server streaming or client streaming? How large do messages have to be to reach peak line rate?

When running your experiments, compile gRPC and your test code both with and without optimization enabled. How much difference does this make in your performance results?

### Deliverables - Part 3

1. **Code**:
   - Source code for the gRPC client and server implementations.
   - Source code for the test cases used to measure marshalling overhead, round-trip time, and bandwidth.

2. **Marshalling Overhead Measurements**:
   - A report detailing the time taken to marshal different types of data (e.g., int, double, string of varying size, complex structure).

3. **Round-Trip Time Measurements**:
   - Measurements of round-trip time for a small message when both client and server are on the same machine.
   - Measurements of round-trip time for a small message when client and server are on different machines.
   - Analysis of the first round trip time compared to subsequent ones.
   - Comparison of round-trip time using gRPC versus the custom reliable communication library.

4. **Bandwidth Measurements**:
   - Measurements of bandwidth when sending large amounts of data using gRPC.
   - Analysis of total achievable bandwidth using server streaming and client streaming.
   - Determination of message sizes required to reach peak line rate.

5. **Optimization Analysis**:
   - Comparison of performance results with and without optimization enabled.

## Machines To Use

For this project, you just need to find two Linux machines/VMs/containers that can speak to one another on a network. If you have access to other machines, that is also fine; for example, most of the cloud service providers give out free credits for educational use. Or, pay for them! They are cheap.

## Handing It In

Deadline for handing in the project: `09/02/2025`, 11:59 PM.  
You should submit a zip file containing all the deliverables for each part of the project. The zip file should be named as `grpc_comm.zip`. The zip file should contain the following directory structure:

```text
grpc_comm/
├── part1/
│   ├── comparison_table.csv
├── part2/
│   ├── code.zip (source code for the reliable communication library and test cases)
│   ├── performance_reliability_measurements.csv
│   ├── reliability_analysis/
│   │   ├── graph1.png
│   │   ├── ...
|   |   └── graphN.png
│   ├── optimization_analysis.csv
├── part3/
│   ├── code.zip (source code for the gRPC client and server implementations and test cases)
│   ├── marshalling_overhead_measurements.csv
│   ├── round_trip_time_measurements.csv
│   ├── bandwidth_measurements.csv
│   ├── optimization_analysis.csv
```

(While the above directory structure is a must; you may submit different extensions for the files, if required.)
