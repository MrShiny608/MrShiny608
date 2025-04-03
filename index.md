---
layout: home
title: My Blog
---

Welcome to my blog! Here are my most recent posts:

<ul>
  {% assign sorted_posts = site.blog | sort: "date" | reverse %}
  {% for post in sorted_posts limit:5 %}
    <li><a href="{{ post.url }}">{{ post.title }}</a> - {{ post.date | date: "%b %-d, %Y" }}</li>
  {% endfor %}
</ul>

