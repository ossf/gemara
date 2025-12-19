---
layout: page
title: Community
nav-title: Community
---

### Join the conversation

- **Slack:** [#gemara](https://openssf.slack.com/archives/C09A9PP765Q) on OpenSSF Slack
- **Meetings:** Bi-weekly on alternate Thursdays - see the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
- **GitHub:** [ossf/gemara](https://github.com/ossf/gemara)

### Meet the Maintainers

{% for maintainer in site.data.maintainers.maintainers %}
- {{ maintainer.name }}, {{ maintainer.organization }} (@{{ maintainer.github }})
{% endfor %}

### Projects and Working Groups 

We are always looking for additional perspective on the Gemara project. Community members and maintainers are involved in several projects, working groups, and initiatives.

* [FINOS Common Cloud Controls Catalog](https://github.com/finos/common-cloud-controls/tree/main/catalogs/core/ccc): FINOS CCC is a collaborative project aiming to develop a unified set of cybersecurity, resiliency, and compliance controls for common services across the major cloud service providers.
* [OpenSSF ORBIT Working Group](https://github.com/ossf/wg-orbit): The development and maintenance of interoperable resource for identification and presentation of security-relevant data. Gemara falls under the ORBIT WG.
* [Open Source Project Security Baseline](https://baseline.openssf.org/): OSPS Baseline is an effort to establish controls that help project maintainers understand security best practices and expectations. 

