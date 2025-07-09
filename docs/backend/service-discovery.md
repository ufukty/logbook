# Service discovery

Design constraints:

- Discovery by **configuration** doesn't scale. Discovery by **self-registration** needs registry service's address to be known at start.

Design:

- `discovery file` contains addresses of `registry service` and `internal gateway` instances

- `discovery file` gets uploaded into `internal gateway` and other service instances by developer (manual process)

- Service instances find addresses of `internal gateway` instances eventually from `discovery file`

- Service instances self-registers themselves by sending a `POST internal/servicereg/register` request to internal gateway with body:

  ```json
  {
    "service": "account",
    "ip": "10.140.0.10",
    "port": 56876
  }
  ```

- `internal gateway` duplicates and sends the `POST /register` request to **all instances** of `registry service` instances

- Services use `internal/web/discovery.ConfigBasedServiceDiscovery` to read `discovery file` and list available instances of the service in need.

Pros:

- The only manual process, file upload only needed when either of `registry service` or `internal gateway` instances change, which occurs less frequently than change of other service instances.

Cons:

- Every instance provision now also requires a file upload. (managable)
- `internal gateway` and `registry service` instances can listen only predefined ports, runtime selection is unallowed. Thus, only one process per address is possible for those.

Improvement opportunities:

- Enable the persistency among `registry service` instances by a hash ring. To remove the need of `internal gateway` to duplicate & forward each registration and recheck request.

```mermaid
sequenceDiagram
autonumber

actor l as developer

create participant d as Registry Service<br>Instance 1<br>IP: U
l->>d: provisions

create participant i as Internal Gateway<br>Instance 1<br>IP: X
l->>i: provisions

note over l: compose discoveryfiles for U, X

l->>i: uploads discoveryfile

create participant s1 as <Service A><br>Instance 1<br>IP: Y
l->>s1: provisions + uploads discoveryfile

rect rgba(128,128,128,0.1)
  note over s1: checks discoveryfile<br>for address of internal
  s1->>+i: POST X/register Y for A
  note over i: checks discoveryfile<br>for address of registry
  i->>+d: POST U/register Y for A
  note over d: stores Y for A in memory
  d->>-i: ok
  i->>-s1: ok
end

create participant s2 as <Service B><br>Instance 1<br>IP: Z
l->>s2: provisions + uploads discoveryfile

rect rgba(128,128,128,0.1)
  note over s2: checks discoveryfile<br>for address of internal
  s2->>+i: POST X/register Z for B
  note over i: checks discoveryfile<br>for address of registry
  i->>+d: POST U/register Z for B
  note over d: stores Z for B in memory
  d->>-i: ok
  i->>-s2: ok
end

rect rgba(128,128,128,0.1)
  s1->>+i: GET X/list/B
  i->>+d: GET U/list/B
  note over d: checks memory for B
  d->>-i: [Z]
  i->>-s1: [Z]
end

s1->>s2: GET Z/foo/bar
```
