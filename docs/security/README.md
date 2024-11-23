# Security

## Methodology

### Equalizer

An equalizer requires a client to spend significant effort to make a request to a service. Which reduces the weakness of server against automated requests.

**Criteria for questionary**

- Computation cost of questionary to client should be at least the same with the cost of user actions to the server
- A Javascript client should be able to solve the questionary as easy as others (or legitimate users are in disadvantage)
- Difficulty should not vary too much, so the completion time.

**Questionary**

There are `m` questions created by `server` and solved by `client`. Flow for a question is like:

1. Server creates question:
   1. Creates a random string in base32: `original`
   1. Hashes original: `hash`
   1. Replaces the first `n` digits of the `original` with spaces: `que`
   1. `[que, hash, n]`: `question`
1. Client tries to solve each question:
   1. hashes the every `32^n` combination on `que`: C
   1. compares `D` and `hash`
   1. answers when `D == hash`

Notes:

- Difficulty increases exponentially with `n`, linearly with `m`.
- Increase the `m` against `n` to reduce deviation in completion time.
- Responder only needs to find answers of a percentage of all questions (~90%).
- Responder should not spend more than 2 times of expected completion time on any question.

**Transparency**

A form utilizes an equalizer may notify users about the usage through friendly text such:

> This form is made computationally difficult for browsers to make one valid submission in order to increase the cost of automated signups for attackers. This technic reduces the need to require users solve puzzles to prove they are not robots. [More info]()

**Function properties**

- Output should be unguessable (uncachable).
- Using ASICs for clients should not provide an advantage over a regular JS client on average device.

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
