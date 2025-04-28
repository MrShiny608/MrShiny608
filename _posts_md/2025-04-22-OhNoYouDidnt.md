---
title: "O(no) You Didn't 😱"
date: 2025-04-22
layout: post
tags: profiling optimisation
categories: optimisation 
excerpt_separator: <!--more-->
---

> A deep dive into why real-world performance often defies Big-O expectations and why context and profiling matter more than theoretical complexity
<!--more-->

[![RSS Feed](https://img.shields.io/badge/RSS-Subscribe-orange?logo=rss&logoColor=white)](https://mrshiny608.github.io/MrShiny608/feed.xml)  [![Check me out on Linkedin](https://img.shields.io/badge/LinkedIn-Profile-0077B5?logo=linkedin&logoColor=white)](https://www.linkedin.com/in/timothybrookes) [![View on GitHub](https://img.shields.io/badge/GitHub-View%20Repo-blue?logo=github)](https://github.com/MrShiny608/code_profiling_playground/tree/master)

## How Can We Make This Faster? 🏎️💨

In one of the many interviews I've been doing lately, I found myself staring down a Leetcode-style question. Now, I know these get mixed reactions - some argue you can't possibly assess a developer's true abilities in 30-60 minutes with a couple of function implementations, while others see the value in how a candidate tackles problems, communicates, and collaborates. Personally? I don’t love them, but I do enjoy the excuse to show off some of the deeper magics...🧙‍♂️

The task? The classic "Two Sum" problem. You're given a list of unique numbers and a target. Find the indices of the two numbers that add up to the target. Easy. I quickly typed out the brute-force solution in Go:

```go
// target := int64(8)
// data := []int64{5, 7, 3, 2}

length := int64(len(data))

for i := int64(0); i < length; i++ {
    complement := target - data[i]

    for j := i + 1; j < length; j++ {
        if data[j] == complement {
            return []int64{i, j}
        }
    }
}

return nil
```

"Okay, nice. How can we make this faster?" the interviewer asked.

*I smiled. The room darkened. Overhead lights flickered and dimmed as if a storm had rolled in. My glasses flared with an unnatural white gleam. Purple tendrils of the low-level deep magics curled up from the floor.* 🌩️

"We can't". 😈

## Big-O Notation 📈

The next few minutes were spent discussing Big-O notation. My implementation? `O(n²)`. The classic answer is to improve it to `O(n)` with a hashmap. But I pointed out - Big-O doesn't measure speed, it measures how *performance degrades* as `n` increases.

Here’s what that looks like:

```go
// target := int64(8)
// data := []int64{5, 7, 3, 2}

hashmap := make(map[int64]int64)
for i, a := range data {
    complement := target - a

    index, ok := hashmap[complement]
    if ok {
        return []int64{index, int64(i)}
    }

    hashmap[a] = int64(i)
}

return nil
```

Let me put it this way:

$$
\begin{aligned}
\text{Assume: } &\quad O_1(n) = c_1 \cdot n \\ &\quad O_2(n^2) = c_2 \cdot n^2 \\
\text{At } n = 4: & \\
c_1 \cdot 4 &\leq c_2 \cdot 16 \\
\Rightarrow \frac{c_1}{c_2} &\leq 4 \\
\end{aligned}
$$

That is, with our tiny input of four elements, our `O(n)` implementation can only afford to have the cost of `O` be four times more expensive than our `O(n²)` implementation before it's actually *slower*. And in the case of the hashmap implementation, `O` is *a lot* slower... 🐢

## Memory Access Times 💾

So what makes the hashmap version slower? It's all about allocations and writes. Every time the CPU writes to memory, it checks the L1 cache, then L2, L3, and finally RAM. The further it goes, the slower it gets:

- L1 cache: ~1ns ⚡
- L2: ~5ns ⚙️
- L3: ~10ns 🧲
- RAM: ~100ns 🐌

The hashmap implementation requires writing up to `n` entries, as well as allocating the backing store for them, potentially multiple times dues to resizes. That's a lot of work for a small input.

## Profiling 📊

As with all performance questions, the only way to answer is: measure. If you're not measuring, you're guessing. I wrote a full post on how to do that right over [here](https://mrshiny608.github.io/MrShiny608/optimisation/2025/04/15/ProfilingCodeWithoutGettingTricked).

For all the graphs below, I deliberately set a worst-case scenario: the target is `-1`, and the data is all unique positive numbers, so no early-outs and no duplicate key advantages.

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Go] Brute Force"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line [63, 134, 242, 390, 581, 809, 1074, 1381, 1726, 2107, 2531, 2987, 3488, 4023, 4604, 5207, 5871, 6565, 7301, 8062, 8914, 9753, 10660, 11553, 12577, 13635, 14728, 15855, 17028, 18254, 19514, 20794, 22166, 23552, 24963, 26396, 27845, 29442, 30935, 32538, 34111, 35846, 37543, 39304, 41143, 44275, 44796, 46693, 48702, 50816, 52743, 54780, 57100, 58952, 61191, 63287, 65617, 67893, 70357, 72770, 75100, 77653, 80053, 82567, 85167, 87938, 90326, 93070, 95600, 98520, 101373, 104334, 107149, 109867, 112871, 116100, 118903, 122060, 125293, 128305, 131344, 134777, 138026, 141309, 144646, 148155, 151435, 155209, 158870, 162494, 165643, 169473, 172726, 176466, 180333, 183986, 188182, 191773, 196087, 199338]
```

</div>
</div>

As expected, the brute-force implementation shows a textbook `O(n²)` curve. At `n = 10`, it takes \~63ns.

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Go] Hashmap"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line [721, 1702, 3475, 3670, 3949, 6920, 7105, 7345, 7604, 7896, 8323, 13674, 13657, 13958, 14204, 14506, 14718, 14962, 15111, 15453, 15728, 16156, 26201, 25976, 26387, 26564, 26745, 26988, 27261, 27576, 27733, 28040, 28081, 28474, 28597, 28963, 29085, 29403, 29657, 30076, 30436, 31440, 30857, 31380, 50576, 50112, 50584, 51210, 50622, 50648, 50920, 52182, 51813, 52149, 52421, 52880, 52816, 52938, 53170, 53644, 53727, 53713, 54282, 54347, 54582, 55168, 55358, 55608, 55638, 55779, 56145, 56663, 56825, 57294, 57441, 57707, 57779, 57911, 59156, 59678, 59483, 59494, 59969, 61112, 61300, 60444, 60593, 61650, 61774, 102351, 102359, 102824, 103689, 103407, 103911, 104128, 104954, 104188, 104214, 104342]
```

</div>
</div>

Now this one’s interesting. 🤔 Ignoring the “steps” for a moment, the general trend is `O(n)` as expected - but at `n = 10` it takes a staggering \~721ns, more than ten times slower than the brute force implementation.

But what are those steps?

They come from bucket allocation. Go’s map implementation starts with an initial bucket size. When it fills up, it grows the map - not by 1, but by a factor. Each resize involves memory allocation *and* copying. That’s expensive. 💸

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Go] Brute Force vs [Go] Hashmap"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line "Brute Force" [63, 134, 242, 390, 581, 809, 1074, 1381, 1726, 2107, 2531, 2987, 3488, 4023, 4604, 5207, 5871, 6565, 7301, 8062, 8914, 9753, 10660, 11553, 12577, 13635, 14728, 15855, 17028, 18254, 19514, 20794, 22166, 23552, 24963, 26396, 27845, 29442, 30935, 32538, 34111, 35846, 37543, 39304, 41143, 44275, 44796, 46693, 48702, 50816, 52743, 54780, 57100, 58952, 61191, 63287, 65617, 67893, 70357, 72770, 75100, 77653, 80053, 82567, 85167, 87938, 90326, 93070, 95600, 98520, 101373, 104334, 107149, 109867, 112871, 116100, 118903, 122060, 125293, 128305, 131344, 134777, 138026, 141309, 144646, 148155, 151435, 155209, 158870, 162494, 165643, 169473, 172726, 176466, 180333, 183986, 188182, 191773, 196087, 199338]

    line "Hashmap" [721, 1702, 3475, 3670, 3949, 6920, 7105, 7345, 7604, 7896, 8323, 13674, 13657, 13958, 14204, 14506, 14718, 14962, 15111, 15453, 15728, 16156, 26201, 25976, 26387, 26564, 26745, 26988, 27261, 27576, 27733, 28040, 28081, 28474, 28597, 28963, 29085, 29403, 29657, 30076, 30436, 31440, 30857, 31380, 50576, 50112, 50584, 51210, 50622, 50648, 50920, 52182, 51813, 52149, 52421, 52880, 52816, 52938, 53170, 53644, 53727, 53713, 54282, 54347, 54582, 55168, 55358, 55608, 55638, 55779, 56145, 56663, 56825, 57294, 57441, 57707, 57779, 57911, 59156, 59678, 59483, 59494, 59969, 61112, 61300, 60444, 60593, 61650, 61774, 102351, 102359, 102824, 103689, 103407, 103911, 104128, 104954, 104188, 104214, 104342]
```

</div>
</div>

Overlaying both implementations, it’s clear: below `~n = 370`, brute force is faster. After `~n = 500`, the hashmap version pulls ahead. 🚀

So back to the interview: *the binary of the universe flickering in and out of existence, the fabric of reality loosening like the seams...* 🪐

"No, we can’t make this faster - not for this use case."

## Context Matters 🎯

Sure, Two Sum is simple. But the lesson applies broadly.

Consider these examples:

- **In-match leaderboard** 🏆: Are you really handling 370+ players at once? Highly unlikely - unless you're running an MMO raid boss from hell. Stick with the `O(n²)` brute-force - it’s simpler, and faster where it matters.
- **Global leaderboard** 🌍: Got thousands or millions of players syncing in real-time? Now you're in scalability territory. This is where `O(n)` starts earning its keep.
- **UI hit detection** 🧩: Ten overlapping widgets on a screen? Brute-force is fine. No need to build a spatial index for your todo app.
- **Ray tracing** 🎥: Hundreds of thousands of rays and geometry? That’s when `O(log n)` acceleration structures save your bacon.
- **Fraud detection** 🔒: Comparing hundreds of transactions per user per second? Better optimize that logic path, or your infra bill will find you.
- **AI search trees** ♟️: Got a 3-move lookahead in chess? Brute force might work. But a 20-move tree? Welcome to exponential growth - better bring pruning and heuristics.

The right choice depends on context, and to prove your assumptions requires measurement. Choosing `O(n)` "because it's faster" is the wrong instinct. It’s not faster - it just has a better *rate of performance decay*.

## Is This Always True? 🤷‍♂️

I wanted to test Python too. It’s not a fair fight - Go is a compiled, low-level beast. Python, even on a good day, is an interpreted langauge who takes strolls in the park, watches the ducks, and honours the UK Tea Alarm. 🫖

But I was curious, how would an interpreted language fare? 🧐

```python
# target: int = 8
# data: List[int] = [5, 7, 3, 2]

length = len(data)
for i, a in enumerate(data):
    complement = target - a

    for j in range(i + 1, length):
        if data[j] == complement:
            return [i, j]

return None
```

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Python] Brute Force vs [Go] Brute Force"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line "Python" [4041, 11507, 19578, 31978, 48040, 66342, 87982, 113226, 141862, 171884, 207907, 244302, 286529, 326551, 376466, 422632, 476458, 535279, 595232, 655895, 718205, 780991, 865129, 926814, 1015479, 1104634, 1197792, 1308087, 1408467, 1515871, 1634337, 1760289, 1889496, 1997579, 2148969, 2285924, 2395276, 2543580, 2709688, 2868030, 2975600, 3129073, 3295737, 3457184, 3651085, 3842209, 3958543, 4200499, 4358340, 4553122, 4706794, 4875898, 5147133, 5340670, 5497817, 5772960, 5923092, 6156185, 6383159, 6589632, 6866112, 7094839, 7268404, 7597067, 7719123, 8096226, 8261885, 8496812, 8819013, 8954966, 9339787, 9601298, 9786468, 10066060, 10409632, 10798281, 10855539, 11176244, 11558334, 11906672, 12226068, 12447802, 12693514, 13166444, 13313960, 13682690, 13962954, 14477178, 14590655, 15156450, 15421691, 15780296, 16162722, 16429044, 16923010, 17025922, 17630922, 17826417, 18364198, 18506079]

    line "Go" [63, 134, 242, 390, 581, 809, 1074, 1381, 1726, 2107, 2531, 2987, 3488, 4023, 4604, 5207, 5871, 6565, 7301, 8062, 8914, 9753, 10660, 11553, 12577, 13635, 14728, 15855, 17028, 18254, 19514, 20794, 22166, 23552, 24963, 26396, 27845, 29442, 30935, 32538, 34111, 35846, 37543, 39304, 41143, 44275, 44796, 46693, 48702, 50816, 52743, 54780, 57100, 58952, 61191, 63287, 65617, 67893, 70357, 72770, 75100, 77653, 80053, 82567, 85167, 87938, 90326, 93070, 95600, 98520, 101373, 104334, 107149, 109867, 112871, 116100, 118903, 122060, 125293, 128305, 131344, 134777, 138026, 141309, 144646, 148155, 151435, 155209, 158870, 162494, 165643, 169473, 172726, 176466, 180333, 183986, 188182, 191773, 196087, 199338]
```

</div>
</div>

As expected, Python's brute force is *slow*, starting at \~4041ns for `n = 10`. I tried variations of the loops: `range`, `enumerate`, `while`... didn't matter. Just running the nested loops, without comparisons, cost \~200ns per entry in the dataset.

So lets check out the hashmap...

```python
# target: int = 8
# data: List[int] = [5, 7, 3, 2]

hashmap = {}
for i, a in enumerate(data):
    complement = target - a

    if complement in hashmap:
        return [hashmap[complement], i]

    hashmap[a] = i

return None
```

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Python] Hashmap vs [Go] Hashmap"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line "Go" [721, 1702, 3475, 3670, 3949, 6920, 7105, 7345, 7604, 7896, 8323, 13674, 13657, 13958, 14204, 14506, 14718, 14962, 15111, 15453, 15728, 16156, 26201, 25976, 26387, 26564, 26745, 26988, 27261, 27576, 27733, 28040, 28081, 28474, 28597, 28963, 29085, 29403, 29657, 30076, 30436, 31440, 30857, 31380, 50576, 50112, 50584, 51210, 50622, 50648, 50920, 52182, 51813, 52149, 52421, 52880, 52816, 52938, 53170, 53644, 53727, 53713, 54282, 54347, 54582, 55168, 55358, 55608, 55638, 55779, 56145, 56663, 56825, 57294, 57441, 57707, 57779, 57911, 59156, 59678, 59483, 59494, 59969, 61112, 61300, 60444, 60593, 61650, 61774, 102351, 102359, 102824, 103689, 103407, 103911, 104128, 104954, 104188, 104214, 104342]

    line "Python" [1604, 2868, 3894, 4939, 6226, 7157, 8220, 9232, 10494, 11453, 12352, 13278, 14282, 15416, 16538, 17575, 18583, 19947, 20873, 21751, 22828, 23763, 24751, 25605, 26548, 27554, 28812, 30133, 31176, 32470, 33739, 34816, 35814, 36874, 39014, 39957, 40821, 41930, 42884, 43967, 44978, 45891, 46912, 47911, 48907, 49896, 50946, 51801, 53038, 54030, 55027, 56388, 57299, 58556, 59767, 60931, 62028, 63223, 64342, 65467, 66680, 67716, 68865, 70270, 71073, 71712, 73012, 73969, 76705, 77639, 78679, 79752, 80757, 81473, 82749, 83532, 84525, 85546, 86822, 87712, 89012, 90185, 90910, 91665, 92591, 93489, 94749, 95804, 96647, 97627, 98373, 100168, 100868, 101932, 102948, 103856, 105052, 105710, 106408, 107770]
```

</div>
</div>

Wait... what? 😳

*Python materializes beside me in the interview room - cloaked in swirling purple vapour, its presence bending the rules of logic and* - okay okay, enough with the dramatisation! 🧙‍♂️

Erm, yes, so, beating all expectations Python has not only achieved performance *on par with Go* but it's also avoided the stepped allocations, at `n = 10` Python took \~1604ns! I checked out some PEPs and the CPython dictionary [source code](https://github.com/python/cpython/blob/main/Objects/dictobject.c), here’s what’s going on:

- `complement in hashmap` and `hashmap[a] = i` are backed by C implementations, these operations don’t stay in Python - they cross the boundary into optimized C code, skipping the interpreter’s overhead entirely. That boundary crossing is expensive in general, but once through, you get raw performance with near-native memory access speeds. 🚀
- Writes to dict are mostly memory-bound - and Python’s bottleneck *is* the CPU. While waiting for those memory accesses to complete, the CPU doesn’t just sit idle. It stays busy interpreting Python bytecode, handling reference counting, checking types, and managing control flow.

Combined, these mean Python is able to perform ridiculously fast, and time it's allocations to fit better within CPU utilisation. If you squint hard, you *can* spot micro-steps, but they're faint. The interpreted nature of Python smooths the curve. 🌊

## So what's the real takeaway? 🎓

Big-O isn't about speed - it's about rate of decay. Profiling beats speculation every time. Complexity tells you *how bad things might get*, but not *when*. That line? You won’t find it on a whiteboard - you find it on a profiler.

- ✅ Measure first. Think later.
- ✅ Choose algorithms based on *real-world contexts*, not theoretical elegance.
- ✅ And maybe cut Python some slack, occasionally it vastly outperforms expectations.

## Closing Notes🧾

### Preallocating Go 🏗️

In Go, we can preallocate enough memory for our hashmap right from the start, eliminating the need for incremental resizes during insertion:

```go
hashmap := make(map[int64]int64, len(data))
```

The rest of the code stays exactly the same. The result? The allocation "steps" become smaller and smoother. There’s still an upfront allocation cost, but now it happens once - early and predictably, which is why we still get steps, but they are much smaller. ✨

The reason I didn't include this in the above is because it added complexity to the post and this post isn't about Go vs Python, it is about `O(n)` vs `O(n²)`

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Go] Hashmap vs Hashmap (Preallocated)"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"

    line "Hashmap" [721, 1702, 3475, 3670, 3949, 6920, 7105, 7345, 7604, 7896, 8323, 13674, 13657, 13958, 14204, 14506, 14718, 14962, 15111, 15453, 15728, 16156, 26201, 25976, 26387, 26564, 26745, 26988, 27261, 27576, 27733, 28040, 28081, 28474, 28597, 28963, 29085, 29403, 29657, 30076, 30436, 31440, 30857, 31380, 50576, 50112, 50584, 51210, 50622, 50648, 50920, 52182, 51813, 52149, 52421, 52880, 52816, 52938, 53170, 53644, 53727, 53713, 54282, 54347, 54582, 55168, 55358, 55608, 55638, 55779, 56145, 56663, 56825, 57294, 57441, 57707, 57779, 57911, 59156, 59678, 59483, 59494, 59969, 61112, 61300, 60444, 60593, 61650, 61774, 102351, 102359, 102824, 103689, 103407, 103911, 104128, 104954, 104188, 104214, 104342]

    line "Hasmap Preallocated" [530, 881, 1443, 1635, 1946, 2504, 2731, 2971, 3181, 3501, 3905, 4690, 4945, 5254, 5433, 5732, 5954, 6208, 6268, 6620, 6947, 7433, 8628, 8867, 9034, 9162, 9578, 9886, 10157, 10403, 10649, 11076, 11123, 11493, 11725, 11971, 12132, 12445, 12352, 12550, 12982, 13443, 13659, 14174, 16459, 16910, 17017, 17127, 17695, 18156, 17906, 18321, 18404, 18655, 19376, 19517, 19942, 19882, 20453, 20321, 20724, 20753, 21550, 21446, 21710, 21802, 22205, 22314, 23013, 22870, 23430, 23596, 23623, 23878, 23939, 24161, 23863, 24548, 24451, 24908, 25817, 25628, 26333, 26564, 26830, 27268, 28065, 28324, 28691, 36944, 37183, 37601, 38030, 38790, 38806, 38902, 39931, 39734, 40630, 41566]
```

</div>
</div>

### Preallocating Python 🐍

As we've seen, choosing the Right Tool For The Job™ 🛠️ can help Python perform significantly better than expected, and preallocating in Go demonstrated a clear performance boost. So, what similar options does Python offer? Unfortunately, after testing dictionary comprehensions and the dictionary `update` and `from_keys` methods, none exhibited behavior resembling preallocation. I even tried caching the `hashmap` dictionary and clearing it out before each call to the two sum function, but it behaved as if `clear` also released the memory, so this didn't help either.

My recommendation here would be to lean towards dictionary comprehensions for pre-populating hashmaps. They showed no adverse impact on performance and offer a clear, idiomatic path that Python developers may further optimise in the future - much like the enhancements made to higher-order functional built-ins such as map, filter, and reduce in the V8 JavaScript engine.

#### Preallocation With Dictionary Comprehension

```python
hashmap = {v: i for i, v in enumerate(data)}

for i, a in enumerate(data):
    complement = target - a

    if complement in hashmap:
        return [hashmap[complement], i]

return None
```

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Python] Hashmap vs Hashmap (Dictionary Comprehension)"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"
    
    line "Hashmap" [1604, 2868, 3894, 4939, 6226, 7157, 8220, 9232, 10494, 11453, 12352, 13278, 14282, 15416, 16538, 17575, 18583, 19947, 20873, 21751, 22828, 23763, 24751, 25605, 26548, 27554, 28812, 30133, 31176, 32470, 33739, 34816, 35814, 36874, 39014, 39957, 40821, 41930, 42884, 43967, 44978, 45891, 46912, 47911, 48907, 49896, 50946, 51801, 53038, 54030, 55027, 56388, 57299, 58556, 59767, 60931, 62028, 63223, 64342, 65467, 66680, 67716, 68865, 70270, 71073, 71712, 73012, 73969, 76705, 77639, 78679, 79752, 80757, 81473, 82749, 83532, 84525, 85546, 86822, 87712, 89012, 90185, 90910, 91665, 92591, 93489, 94749, 95804, 96647, 97627, 98373, 100168, 100868, 101932, 102948, 103856, 105052, 105710, 106408, 107770]

    line "Dictionary Comprehension" [1938, 3051, 4254, 5006, 5976, 6826, 7891, 9018, 9925, 10898, 11838, 12759, 13739, 14852, 15950, 17097, 18310, 18721, 19742, 20671, 21627, 22583, 23425, 24395, 25529, 26416, 27637, 28772, 30095, 31354, 32667, 33882, 35242, 36610, 37843, 38570, 39199, 40453, 42311, 42359, 43714, 44785, 45686, 47046, 47672, 48870, 50122, 51002, 52105, 53317, 54317, 56330, 57076, 57920, 59305, 60630, 62237, 62989, 64048, 65568, 67117, 67860, 70415, 72513, 71838, 73414, 74472, 75834, 75851, 76789, 77449, 78643, 79841, 80958, 81530, 82938, 84493, 86184, 86857, 87424, 88614, 91008, 91164, 91900, 92593, 94050, 95496, 96375, 97004, 98424, 99232, 100158, 101397, 102220, 103391, 104470, 108150, 107277, 108111, 109321]
```

</div>
</div>

As we see there is no significant performance difference, nor is the dictionary comprehension consistently faster.

#### Preallocation With Update

```python
hashmap = {}
hashmap.update(enumerate_to_dict_update(data))

for i, a in enumerate(data):
    complement = target - a

    if complement in hashmap:
        return [hashmap[complement], i]

return None
```

In this case we have to build a custom generator to convert `enumerate` to the correct form for `update`

```python
def enumerate_to_dict_update(data: List[int]) -> Generator[Tuple[int, int], None, None]:
    for i, a in enumerate(data):
        yield a, i
```

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Python] Hashmap vs Hashmap (Update)"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"
    
    line "Hashmap" [1604, 2868, 3894, 4939, 6226, 7157, 8220, 9232, 10494, 11453, 12352, 13278, 14282, 15416, 16538, 17575, 18583, 19947, 20873, 21751, 22828, 23763, 24751, 25605, 26548, 27554, 28812, 30133, 31176, 32470, 33739, 34816, 35814, 36874, 39014, 39957, 40821, 41930, 42884, 43967, 44978, 45891, 46912, 47911, 48907, 49896, 50946, 51801, 53038, 54030, 55027, 56388, 57299, 58556, 59767, 60931, 62028, 63223, 64342, 65467, 66680, 67716, 68865, 70270, 71073, 71712, 73012, 73969, 76705, 77639, 78679, 79752, 80757, 81473, 82749, 83532, 84525, 85546, 86822, 87712, 89012, 90185, 90910, 91665, 92591, 93489, 94749, 95804, 96647, 97627, 98373, 100168, 100868, 101932, 102948, 103856, 105052, 105710, 106408, 107770]

    line "Update" [2562, 4247, 9610, 7271, 8684, 10164, 11735, 13310, 14828, 16312, 17743, 19307, 20780, 22434, 23971, 25554, 27500, 28555, 30099, 31502, 32818, 34334, 35929, 37327, 38642, 40552, 42215, 44255, 45810, 47696, 49872, 51739, 53500, 55496, 56131, 58430, 59673, 61409, 63012, 64639, 66522, 68110, 69894, 71474, 72976, 75067, 76511, 78108, 79860, 81460, 83268, 85035, 87274, 88785, 90641, 92442, 94350, 95817, 97802, 99435, 101233, 103387, 105106, 107149, 109045, 110989, 112777, 114953, 114928, 116565, 118166, 119714, 121379, 123507, 124957, 126542, 128203, 129973, 131664, 133142, 135345, 136750, 138332, 139739, 142005, 144093, 144721, 147004, 148765, 150325, 151965, 153837, 155438, 156936, 158705, 160554, 161706, 163447, 165372, 167504]
```

</div>
</div>

Here we see that `update` is slower than using our original implementation, which is a surprise as at the most it should be doing the same work...

#### Preallocation With From Keys

```python
hashmap = dict.fromkeys(data)

for i, a in enumerate(data):
    complement = target - a

    j = hashmap.get(complement)
    if j is not None:
        return [j, i]

    hashmap[a] = i

return None
```

<div class="mermaid-grid">
<div class="xlarge-inline-card">

```mermaid
xychart-beta
    title "[Python] Hashmap vs Hashmap (From Keys)"
    x-axis "iterations" 10 --> 1000
    y-axis "nanoseconds"
    
    line "Hashmap" [1604, 2868, 3894, 4939, 6226, 7157, 8220, 9232, 10494, 11453, 12352, 13278, 14282, 15416, 16538, 17575, 18583, 19947, 20873, 21751, 22828, 23763, 24751, 25605, 26548, 27554, 28812, 30133, 31176, 32470, 33739, 34816, 35814, 36874, 39014, 39957, 40821, 41930, 42884, 43967, 44978, 45891, 46912, 47911, 48907, 49896, 50946, 51801, 53038, 54030, 55027, 56388, 57299, 58556, 59767, 60931, 62028, 63223, 64342, 65467, 66680, 67716, 68865, 70270, 71073, 71712, 73012, 73969, 76705, 77639, 78679, 79752, 80757, 81473, 82749, 83532, 84525, 85546, 86822, 87712, 89012, 90185, 90910, 91665, 92591, 93489, 94749, 95804, 96647, 97627, 98373, 100168, 100868, 101932, 102948, 103856, 105052, 105710, 106408, 107770]

    line "From Keys" [3127, 5409, 7229, 9962, 11890, 13690, 15966, 18027, 19891, 22376, 24348, 26813, 28149, 30992, 33092, 35451, 38163, 39006, 40907, 42259, 45271, 47526, 49836, 50105, 52534, 55413, 56097, 60695, 62795, 66580, 68747, 70142, 74579, 75620, 74950, 78104, 79047, 81779, 83628, 85323, 87589, 89925, 92899, 94631, 95826, 100188, 98674, 102842, 103344, 106192, 110640, 111580, 114821, 114733, 118284, 119124, 121433, 126523, 130207, 132764, 130689, 136652, 140162, 141920, 143675, 143606, 151480, 150071, 150396, 151526, 151142, 154332, 156618, 158015, 161981, 165078, 170020, 170785, 171137, 171447, 177672, 175624, 182194, 177501, 181392, 181701, 185824, 188183, 188532, 191768, 195430, 199384, 197024, 202958, 202214, 210275, 150977, 140232, 141372, 142729]
```

</div>
</div>

Here we see that `from_keys` is also slower than our original implementation, even including the dramatic reduction around a data size of 1000

---

[![RSS Feed](https://img.shields.io/badge/RSS-Subscribe-orange?logo=rss&logoColor=white)](https://mrshiny608.github.io/MrShiny608/feed.xml)  [![Check me out on Linkedin](https://img.shields.io/badge/LinkedIn-Profile-0077B5?logo=linkedin&logoColor=white)](https://www.linkedin.com/in/timothybrookes) [![View on GitHub](https://img.shields.io/badge/GitHub-View%20Repo-blue?logo=github)](https://github.com/MrShiny608/code_profiling_playground/tree/master)
