---
title: 💎 Quality of Guarantees
date: 2025-05-13
layout: post
tags: dev_culture consistency conventions guarantees
categories: commentary
excerpt_separator: <!--more-->
---

> What if the rock-solid foundations your system relies on... aren’t?
<!--more-->

## 🧠 Foundations Built on Sand?

Many software engineers like to think they work in a very "pure" field - much like in math, where 1 + 1 always ([okay, normally](https://medium.com/%40mmajormoss/explaining-as-a-mathematically-disinclined-individual-why-1-1-does-not-equal-2-9deaac7d4c09)) equals 2, we assume a given operation with the same inputs always yields the same output. But does it? In reality, software is built atop countless conventions and guarantees - some strong, some... not so much.

## 📜 Conventions

There are a lot of things we take for granted: Unix time starts on Jan 1st 1970; a zero exit code means success; `index.html` is the default page for a website. These are conventions. They're not universal truths - they exist because we, as a software community, decided to standardise them.

Some conventions have sensible origins. Zero as false comes from checking if the electrical charge is greater than the required threshold. `\n` terminates a line, and Windows' `\r` returns the cursor to the start of the line thanks to the mechanical process of typewriters. Arrays start at 0 due to pointer arithmetic - see [Dijkstra's explanation](https://www.cs.utexas.edu/~EWD/transcriptions/EWD08xx/EWD831.html). Others are more arbitrary: HTTP uses port 80 because [Tim Berners-Lee picked it](https://www.w3.org/Protocols/HTTP/AsImplemented.html).

And yet, these conventions underpin everything from our shell scripts to our distributed services.

## 🛡️ Guarantees

Beyond conventions, we rely heavily on guarantees. Memory we set stays set - unless a cosmic ray flips a bit. UUIDv4 won't collide - probably. Disk writes are durable after `fsync()` - usually.

Some guarantees sound strong but aren’t:

* ⏱️ Time always moves forward? NTP corrections can rewind the clock.
* 🌍 The internet always routes traffic correctly? [BGP misconfigurations](https://www.ripe.net/about-us/news/youtube-hijacking-a-ripe-ncc-ris-case-study/) can misroute entire countries.
* ➗ Floating point math is consistent? IEEE 754 compliance varies - and optimisations may reorder operations, making `a + b + c` not equal to `b + a + c` due to [non-associativity in floating point numbers](https://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html).

## 📈 An Issue of Scale

Yes, the internet is unreliable. Any request can fail. Networks partition. Clocks drift. But once you scale to thousands of requests per second, those rare edge cases become frequent visitors.

That shiny "five nines" SLA (99.999% availability) sounds robust - but that still allows over **5 minutes** of downtime per year.  In practice, issues tend to cluster rather than distribute evenly - we notice some and fix them, but a steady trickle of subtle, overlooked blips still slip through. These unnoticed glitches can become landmines when they collide with systems that assumed certain behaviours were guaranteed.

Once you start hitting these edges, the cracks in your assumptions widen. What you thought was "guaranteed" becomes "likely". What was "rare" becomes "frequent".

## 🧨 When Guarantees Aren't

Some guarantees look solid until you realise they’re balanced on assumptions that aren’t guaranteed at all. For example, a pattern I see gaining popularity is the Event Outbox pattern - designed to solve the issue of writing data and emitting an event in a single atomic operation.

```python
database.update("user1234", "some_data")
messageQueue.push("Hey, I updated user1234")
```

Looks good? Not quite. If the system crashes in between those lines - or the coroutine thread is torn down - you’ll lose the message. Runtimes like Node.js or collaborative concurrency libraries like asyncio in Python are particularly exposed here, but no environment is immune. A SIGKILL or OOM can cut off your process between those two lines.

Enter the Event Outbox: write the message to a table in the same transaction as the data update. A background service reads from this table and pushes to the queue, marking it as sent only once acknowledged.

This solves the crash-in-the-middle issue, but it doesn't fix more critical vulnerabilities:

* ❓ What if a developer forgets to add the message to the outbox?
* 🔁 What if a migration script corrects malformed data but skips messaging?
* 🧑‍💻 What about manual DB edits? (Yes, I’ve seen cultures where engineers run direct SQL against prod to "fix" bugs.)

I’ve worked with companies where entire teams - as much as 10% of the engineering org - were dedicated solely to repairing data inconsistencies: patching over bugs, correcting silent human mistakes, and recovering from missing or misfired events.

I brought this up in a recent interview. The tech director dismissed the problem with a shrug - "you can't solve everything." And sure, not everything needs solving. But here, you actually *can* harden a fragile, human-enforced guarantee into a deterministic computing one.

Rather than watching the outbox table with a CDC tool and hoping developers remember to populate it, monitor all database changes directly. Let the event stream be derived from the source of truth itself - the write-ahead log or a logical replication stream. With that in place, no manual edit, silent script, or forgotten migration can alter data without triggering an event.

## ✅ Wrapping Up

Most of what we rely on in software isn’t as stable as we think. Some things are conventions with strong cultural momentum, others are guarantees with asterisks and fine print. Sometimes we build systems on rock; other times, we’re skating over ice and hoping it doesn’t crack.

Knowing the difference - and designing with eyes open - is often what separates a system that wakes up the on-call team at 3 a.m. from one that keeps on ticking. The job isn’t to eliminate every risk. It’s to understand which ones you’ve taken - and to recognise when it’s time to reassess, adapt, or reinforce those assumptions before they turn brittle.
