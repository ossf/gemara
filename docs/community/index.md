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

### Projects and Initiatives on our Radar

We're glad you're here! We are always looking for additional perspective on the Gemara project and how the scope can be extended. Community members and maintainers are involved in several projects, working groups, and initiatives.


* [**Initiative for Controls Catalog Refresh**](https://github.com/cncf/toc/issues/1910) led by Gemara maintainer Jenn Power, an effort to review and update the _existing_ controls within the Cloud Native Security Controls Catalog. The initiative supports automated assessment of the security controls.   
  * Find the Slack Channel [here](https://cloud-native.slack.com/archives/C09TLL22PK9)
* Explore the [FINOS Common Cloud Controls Catalog](https://github.com/finos/common-cloud-controls/tree/main/catalogs/core/ccc) expressed in Gemara
* [OpenSSF ORBIT Working Group](https://github.com/ossf/wg-orbit): development and maintenance of interoperable resource for identification and presentation of security-relevant data. 
* [Open Source Project Security Baseline](https://baseline.openssf.org/): OSPS Baseline is an effort to establish controls that help project maintainers understand security best practices and expectations. 

