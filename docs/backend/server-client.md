# Server — Client

> | Kind       | Item       | Description              |
> | ---------- | ---------- | ------------------------ |
> | Service    | account    | Account                  |
> | Service    | auth       | Auth                     |
> | Service    | bookmark   | Bookmark                 |
> | Service    | history    | Task History             |
> | Service    | invitation | Task Invitation          |
> | Service    | privilege  | Task Privileges          |
> | Service    | push       | Push Notification Server |
> | Service    | task       | Task                     |
> | Service    | taskLink   | Task Link                |
> | Service    | timemngr   | Time Manager             |
> | Service    | view       | Document View            |
> | User       | collab     | Worker                   |
> | User       | owner      | Task owner               |
> | Web server | spa        | SPA                      |

## Registration

```mermaid
sequenceDiagram
autonumber

owner ->>+ account: create account
account ->> task: create root task
task ->> account: return task_id
account ->>- bookmark: create "all" bookmark
```

## Login

```mermaid
sequenceDiagram
autonumber

owner ->>+ auth: credentials
auth ->>- owner: authorization token
```

## Enter user land & welcome page

```mermaid
sequenceDiagram
autonumber

activate owner
owner ->> spa: READ page
spa ->> owner: HTML + JS
par
  owner ->> bookmark: READ bookmarks
and
  owner ->> timemngr: READ reports
end
deactivate owner
```

## Load document from `all` bookmark

```mermaid
sequenceDiagram
autonumber

owner ->>+ view: READ hierarchical_view_placement(root_task_id: bookmark ref, range: [start, end))
view ->> view: check cache
view ->> taskLink: run DFS, with child passing when:<br>degree + passed_nodes < starting_index<br>and terminate growing view when:<br>passed_nodes > ending_index
view ->> owner: placement array for range [start, end)
owner ->> push: create connection for future incoming data
```

## Create task

```mermaid
sequenceDiagram
autonumber

owner ->> task: create_task(user_id, root_task_id, super_task_id, task_details)
task ->> privilege: GET roles(user_id, super_task_id)
task ->> task: CREATE task
task ->> taskLink: CREATE link between task_id, super_task_id
task ->> taskLink: COMPUTE degree, depth, index of task_id
task ->> taskLink: RECOMPUTE (RECURSIVE to super-task)<br>degree, depth, index of super_task_id
```

## Recomputing properties of a task link

```mermaid
sequenceDiagram
autonumber

taskLink ->> taskLink: recompute
taskLink ->> push: INVALIDATE BROWSER CACHE for task_id
taskLink ->> push: INVALIDATE BROWSER CACHE for placement_array
```

## Invite user as `maintainer` over `task_id`

```mermaid
sequenceDiagram
autonumber

owner ->> privilege: CREATE invitation()
privilege ->> privilege: check roles (authorization)
privilege ->> push: push event
```

## Accept invitation

```mermaid
sequenceDiagram
autonumber

push ->> collab: notification: "new invitation"
collab ->> privilege: ds
```

## Unassign user as `maintainer`

```mermaid
sequenceDiagram
autonumber

owner ->> privilege: ds
```

## task details

```mermaid
sequenceDiagram
autonumber

owner ->> taskLink: GET degree, depth (for user)
owner ->> history: GET history (inc. subtasks optional)
owner ->> privilege: GET roles (inc. inherited)
```
