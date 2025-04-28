---
title: 🧪 Profiling Code Without Getting Tricked
date: 2025-04-15
layout: post
tags: profiling optimisation
categories: optimisation 
excerpt_separator: <!--more-->
---

> A deeper dive into profiling code for optimisation - and how your measurements might be lying to you
<!--more-->

[![Check me out on Linkedin](https://img.shields.io/badge/LinkedIn-0077B5?logo=linkedin&logoColor=white)](https://www.linkedin.com/in/timothybrookes) [![View on GitHub](https://img.shields.io/badge/GitHub-View%20Repo-blue?logo=github)](https://github.com/MrShiny608/code_profiling_playground/tree/master)

## 🧠 Guessing Isn't Good Enough

So you've got a system that's performing poorly. You've used system profiling tools and observability to figure out that a lot of time is spent executing a function named `DoSomething()`. You take a look and find:

```go
func DoSomething(data SomeData[]) {
  for _, entry := range data {
    DoItSlowly(entry)
  }
}
```

Simple fix! You think, I'll just optimise it:

```go
func DoSomething(data SomeData[]) {
  for _, entry := range data {
    DoItFast(entry)
  }
}
```

Done, right? 🚀

Unfortunately, this is where many engineers go wrong. They speculate about the problem and implement a fix without ever measuring whether it actually helped. If you're not measuring, you're guessing.

There’s nothing wrong with forming hypotheses - that’s how good science works. But the hypothesis isn't the answer. You’ve got to test it.

So how do we *actually* measure our optimisations?

## ⏱️ Why Measuring Is Surprisingly Hard

Measuring the performance of code accurately is actually... *really* hard. System profiling tools tend to poll frequently and sample where your code is spending time, over many samples you get a statistical picture. This observability helps, but it doesn't let you measure how long a single call took.

### 🕰️ Not All Timers Are Created Equal

Even if your CPU is running at just 1GHz (i.e., one cycle every \~1ns), you can't measure time that accurately. Tools like [`QueryPerformanceCounter`](https://learn.microsoft.com/en-us/windows/win32/api/profileapi/nf-profileapi-queryperformancecounter) (Windows) or [`clock_gettime(CLOCK_MONOTONIC)`](https://linux.die.net/man/3/clock_gettime) (Linux) offer microsecond *resolution*, but their *accuracy* is a different story.

```go
type Callable func() (result any)

func Profile(work Callable, data []SomeData) {
  // 🚫 Not accurate at all - don't do this
  start := time.Now()
  work(data)
  duration := time.Since(start)
  fmt.Printf("Took %dns", duration.Nanoseconds())
}
```

These clocks are too unstable to rely on for short durations. For example:

- `QueryPerformanceCounter` (Windows) typically has a resolution of 1 microsecond, but accuracy drifts under load and may be updated only every \~0.5 to 1µs depending on the platform.
- `clock_gettime(CLOCK_MONOTONIC)` (Linux) can report nanosecond resolution, but real-world jitter is usually in the 5–50µs range due to system scheduling and kernel overhead.
- `time.Now()` in Go (like `time.time()` in Python or `clock()` in C) wraps platform APIs and often only updates every \~16µs on many systems - making it nearly useless for microbenchmarks. This limitation is common across many high-level languages, as they rely on the same underlying system clocks.

Resolution != accuracy. Even if a clock *can* report values to the micro- or nanosecond, that doesn’t mean it reflects reality to that precision 😅.

### 🔁 Beat the Noise With Repetition

So we run the function many times and divide the total duration by the number of runs. Easy fix? Kind of... but now we introduce (or rather, exacerbate) CPU noise - the unavoidable fluctuations caused by everything else your system is doing in the background.

Your CPU isn’t just running your function - it’s juggling threads, handling interrupts, managing OS processes, executing instructions out of order, and possibly adjusting clock speed on the fly.

```go
type Callable func() (result any)

func Profile(work Callable, data []SomeData, iterations int64) {
  // ⚠️ Not accurate unless iterations is large enough
  start := time.Now()

  for i := range iterations {
    work(data)
  }

  duration := time.Since(start)
  fmt.Printf("Took %dns", duration.Nanoseconds()/iterations)
}
```

Now we're averaging background noise too 😬. By running the function many times, we aren't eliminating the noise entirely, but we are smoothing it out. The more iterations we include, the more the random jitter from other processes, OS scheduling, or CPU frequency scaling gets averaged out. This helps us get closer to a *comparable* true cost of the code under test. And since all of our test cases experience this noise in a similar way, we expect that by averaging enough runs we reduce the impact of that noise on our comparisons.

### 📊 Bigger Data, Better Insights

To get clean results, you need to go a little further. Most functions slow down as input grows (Big-O complexity), so you want to test across a range of sizes and durations. This both shows how the profile changes with input size, and gives us a chance to identify outliers that might be skewing the data.

```go
type Callable func() (result any)

type Test struct {
  Work Callable
  Data []SomeData
}

func Profile(tests []Test, testDuration time.Duration) {
  for _, test := range tests {
    iterations := int64(0)
    start := time.Now()

    for time.Since(start) < testDuration {
      test.Work(test.Data)
      iterations++
    }

    sampleDuration := time.Since(start)
    fmt.Printf("N=%d took %dns", len(test.Data), sampleDuration.Nanoseconds()/iterations)
  }
}
```

Using a fixed time budget (e.g. 5 minutes per test) ensures each case runs long enough to average out random noise - without having to manually scale the number of iterations for each input size. This approach adapts naturally as datasets grow, while keeping comparisons consistent.

<div class="mermaid-grid">
<div class="xlarge-inline-card">

<div class="mermaid">
xychart-beta
    title "An O(N^2) Implementation"
    x-axis "iterations" 1000 --> 100000
    y-axis "nanoseconds" 

    line [425651, 1985600, 4347336, 7906971, 10894523, 14459480, 18678951, 23193010, 28453348, 34481522, 40996807, 47872225, 56097571, 63806369, 72675124, 82017625, 92678987, 102791013, 114891825, 125697768, 138634671, 153030973, 167179185, 182183164, 195571815, 212657384, 228662245, 245897563, 261479811, 282253003, 299118150, 320059977, 337099163, 357461883, 381802923, 402583671, 426610878, 448413754, 474444551, 497264172, 522226729, 552674362, 565615488, 598556164, 702608664, 659584697, 677638335, 707711740, 741237462, 773367412, 801404404, 833525724, 862336427, 897239857, 935762962, 970179049, 1005690864, 1036699782, 1069166854, 1103224196, 1144155664, 1172987539, 1225781451, 1263161411, 1291457825, 1333661107, 1379015278, 1413041211, 1463382971, 1494415157, 1547178753, 1577455793, 1636571323, 1674025246, 1751401643, 1776123539, 1800307252, 1868881609, 1901483083, 1958551757, 1973700230, 2054877770, 2129828432, 2158190476, 2223453133, 2246020330, 2307699530, 2379102207, 2435325597, 2478058857, 2505372806, 2561332689, 2664061506, 2685866500, 2744224489, 2810060982, 2873942560, 2923211369, 2989222897, 3074190067]   
</div>

</div>
</div>

When you graph the results...👀 One of those values looks very sus. Graphs FTW. Because our human brains are actually really good at spotting outliers, we can choose to ignore that sample, re-run the test, or at least treat it with suspicion when comparing against other implementations.

## 🚧 Spotting the Lies in Your Results

Even once you can measure consistently, there are still traps waiting to trip you up...

### 🧊 When Memory Placement Betrays You

Just like data, code is loaded from memory into caches. If your function crosses cache lines or gets evicted and reloaded, performance goes boom 💥. Cache misses can cost hundreds of cycles.

Profiling in situ sounds nice - it mimics the real production environment. But the danger is overfitting: your optimisation might be too specific to the current memory layout, code arrangement, or cache alignment. Even a small unrelated change elsewhere in the codebase could shift things around enough to invalidate your carefully tuned gains.

If you’ve ever seen an optimisation mysteriously stop working weeks later, this might be why.

✅ Tip: Profile in a minimal, isolated environment. Control everything you can - ideally, compile each test separately to minimise interference from unrelated code.

### ✂️ Your Compiler is Smarter Than You Think

Sometimes your profiled function seems incredibly fast - suspiciously fast. That’s often because the compiler has optimised it away entirely, removing the work because it noticed the result isn’t used anywhere.

To prevent this, make sure your function returns a value, and do something with it. Don't let the compiler skip the work.

```go
func work(data []SomeData) (result any) {
  result = 0
  for _, entry := range data {
    // The original optimisation
    DoItFast(entry)
    
    // Some work to persist the loop in the isolated test
    result++
  }

  // Make sure to return it so the compiler doesn't
  // remove the variable entirely
  return result
}
```

As long as you're consistent between your tests, this won’t skew results.

### 🧠 Your Processor is Smarter Than You Think

Modern CPUs don’t just sit and execute your code - they learn. As your program runs, the processor dynamically adjusts how it executes instructions to improve performance:

- 🔄 Branch prediction gets more accurate as the CPU observes actual branching patterns.
- 🧠 Instruction caching and micro-op fusion kick in to reduce pipeline stalls.
- 🚀 Out-of-order execution improves throughput as the CPU figures out which operations it can run in parallel.
- 📈 Speculative execution tries to guess and execute future instructions before they’re needed.

All of these things can make your program faster over time - even without code changes. But they also make profiling harder, because early runs can behave differently from later ones.

### 🔥 Interpreters Need a Warm-Up Lap

In JIT languages like Python (PyPy), Java (JVM), or JavaScript (V8), runtime optimisers kick in after observing hot code - often based on a combination of call frequency, loop iteration counts, type stability, and inline cache feedback.

So if you're benchmarking a function, it's critical to warm it up before measuring it.

```python
for _ in range(10000):
    do_work()
```

This gives the runtime time to identify it as hot and optimise the machine code it runs. Just make sure to do this before you start the timer 😉.

### 🧹 Garbage Collectors Get Better With Age

Garbage collectors (GCs) often perform better over time within a single process - as the runtime gathers statistics, warms up generational pools, or optimises allocation strategies. This means your later tests might unfairly benefit from a better tuned GC state, making them look faster even if they aren't.

## ✅ Wrapping Up

Profiling isn't just about slapping a timer around your function and calling it a day. It takes care, statistical thinking, and awareness of the many layers that can trick you - from OS jitter to sneaky compiler tricks, to caching and warm-up quirks.

If you want meaningful, repeatable results:

- 🕰️ Use multiple iterations to negate poor clock accuracy, and use enough iterations to reduce the effect of short-term noise.
- 🧼 Isolate your code to reduce interference from other application code and runtime state.
- 🔁 Run each test configuration in its own minimal binary to eliminate bias from cold vs warm runtime and hardware optimisations.
- 🔄 Execute each binary multiple times to smooth out mid-term variability.
- 🔀 Interleave similar configurations from different tests to reduce bias from longer lived variability.
- 🧠 Apply human judgment when reviewing outliers - graphs and intuition can spot things stats miss.
- 🔥 Warm up your JITs to measure code that’s actually running how it would in production.
- 🤖 Make sure you're really measuring what you think you are - no optimised-away logic or skewed setup costs.

---

[![Check me out on Linkedin](https://img.shields.io/badge/LinkedIn-0077B5?logo=linkedin&logoColor=white)](https://www.linkedin.com/in/timothybrookes) [![View on GitHub](https://img.shields.io/badge/GitHub-View%20Repo-blue?logo=github)](https://github.com/MrShiny608/code_profiling_playground/tree/master)