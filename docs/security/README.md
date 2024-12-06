# Security

## Methodology

### Equalizer

An equalizer requires a client to spend significant effort to make a request to a service. Which reduces the weakness of server against automated requests.

**Criteria**

- **Effectiveness** - Solving the challenge must be more difficult than serving to an anonymous user
- **Frictionless** - User doesn't need to interact to pass the challenge.
- **Equality** - A JavaScript client should not have disadvantage against an attacker using 3rd-party client written in another language using an ASIC
- **Predictability** - Completion time should not vary too much between runs with same parameters
- **Adjustability** - The difficulty needs to be adjusted with minimal increments based on the volume of invalid requests originating from IPs/blocks.

**Design**

| Variable | Description          | Default |
| -------- | -------------------- | ------- |
| CPB      | Challenges per batch | 100     |
| IL       | Input length         | 20      |
| ML       | Mask length          | 3       |

Create

```python
alphabet = "012...ABC...abc..."

def CreateChallenge(difficulty: number):
  original = randomstring(ML,      alphabet[:difficulty]) +
             randomstring(IL - ML, alphabet)
  hashed = hashfunction(original)
  masked = hashed[ML:]
  return (original, masked, hashed)

def CreateBatch(difficulty: number):
  return [CreateChallenge(difficulty) for _ in range(CPB)]
```

Solve

```python
alphabet = "012...ABC...abc..."

def SolveChallenge(difficulty: number, masked: string, hashed: string):
  combination = combinate(ML, alphabet[:difficulty])
  for {
    if hashfunction(combination + masked) == hashed {
      return combination
    }
    if !combination.iterate() {
      return nil
    }
  }

def SolveBatch(difficulty: number, challenges: [](masked, hashed)):
  cs = {}
  for masked, hashed in challenges:
    combination = SolveChallenge(difficulty, masked, hashed)
    if combination == NIL {
      return Error("failed for {combination}")
    }
    cs[hashed] = combination
  return cs
```

Validate

```python
def ValidateBatch(batchid: number, combinations: [](hashed, combination)):
  challenges, ok = store[batchid]
  if !ok:
    return Error("invalid batch id")
  if len(challenges) != len(combinations):
    return Error("wrong number of answers")
  if !compare_hashed_values(challenges, combinations):
    return Error("wrong questions")
  if compare_combinations(challenges, combinations):
    return Error("wrong answers")
  return nil
```

- **Challenges per batch**  
  It is picked as such to make sure the batch gets balanced selection of quick and slow challenges. Which the speed of solution determined by how late the solution is placed in the set of possible values. The algorithm doesn't contain logic to make sure equal number of early/late solutions by intention to not give attackers means to guess the direction for trials for sequent challenges.
- **Hash function**  
  outputs text that is in transport safe encoding. the result should be unguessable.
- **Alphabet limiting**  
  is applied only to the generation of masked part, to reduce weakness against rainbow tables.

**Review**

Increasing difficulty

| Growth      | Parameter                       | Example terms               |
| ----------- | ------------------------------- | --------------------------- |
| exponential | # of masked characters          | `ML^2`, `ML^3`, `ML^4`      |
| polynomial  | # of characters in the alphabet | `x^3`, `(x+1)^3`, `(x+2)^3` |
| linear      | # of challenges per batch       | `CPB`, `2CPB`, `3CPB`       |

> Increase the `CPB` against rest to reduce deviation in completion time.

> Increasing `CPB` over `ML` reduces the disparity between difficulties for the server to create challenges and the client to solve them.

**Transparency**

A form utilizes an equalizer may notify users about the usage through friendly text such:

> This form is made computationally difficult for browsers to make one valid submission in order to increase the cost of automated signups for attackers. This technic reduces the need to require users solve puzzles to prove they are not robots. [More info]()

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

Failed equalizer challenges in last hour, day, week etc.

**Setting rules**

When possible, generalize black list rules to limit their growth until degrees deciding a request become too expensive.

Use information to set black list rules:

- IP addresses

Consider using white listing:

- Emails with `.edu`, `.gov` tlds.
- Reputable email providers `gmail.com`, `outlook.com` (?)

**Penalties**

- Exponentially increase the equalizer challenge difficulty on marked IP blocks.

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
