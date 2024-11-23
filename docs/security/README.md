# Security

## Measures

### Registration

**Concerns**

Failure on blocking automated account creation attempts may provide basis to:

- Abuse of free resources
- Denial of service
- Decreasing trust on userbase authenticty
- Manipulation of statistics on content popularity

**Available information**

There are not much information available throughout the registration to mark attackers:

- IP address
- Profile (name, surname, birthday, country)
- Contact (email, phone)

**Attacker classification**

Failed equalizer challanges in last hour, day, week etc.

**Setting rules**

When possible, generalize black list rules to limit their growth until degrees deciding a request become too expensive.

Use information to set black list rules:

- IP addresses

Consider using white listing:

- Emails with `.edu`, `.gov` tlds.
- Reputable email providers `gmail.com`, `outlook.com` (?)

**Penalties**

- Exponentially increase the equalizer challange difficulty on marked IP blocks.

**Consequences**

Using shared IPs for registration is not supported.

## Backend

### Low level blocking

Application server maintains a log file which is also watched by fail2ban.

- The log file contains the failed attempts to perform API requests.

### User level rate limiting

### Authorization

#### PDP

The design shapes a **centralized** & **stateless** Policy Decision Point (PDP).

**Design constraints**

- **Verifiability**: Gather decision rules to one place to ease verifying implementation against requirements
- **Complex access control logic**: Due to collaboration and hierarchical groups
- **Distributed data**: Involved data in the decision process is too big and varying to gather in one place

![](pdp.png)
