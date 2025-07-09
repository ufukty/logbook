| Kind       | Item       | Description              |
| ---------- | ---------- | ------------------------ |
| Service    | account    | Account                  |
| Service    | auth       | Auth                     |
| Service    | bookmark   | Bookmark                 |
| Service    | history    | Task History             |
| Service    | invitation | Task Invitation          |
| Service    | privilege  | Task Privileges          |
| Service    | push       | Push Notification Server |
| Service    | task       | Task                     |
| Service    | taskLink   | Task Link                |
| Service    | timemngr   | Time Manager             |
| Service    | view       | Document View            |
| User       | collab     | Worker                   |
| User       | owner      | Task owner               |
| Web server | spa        | SPA                      |

## Registration

- owner → account: create account
- account → task: create root task
- task → account: return task_id
- account → bookmark: create "all" bookmark

## Login

- owner → auth: credentials
- auth → owner: authorization token

## Enter user land & welcome page

- owner → spa: READ page
- spa → owner: HTML + JS
- paralle:
  - owner → bookmark: READ bookmarks
  - owner → timemngr: READ reports

## Load document from "all" bookmark

- owner → view: READ `hierarchical_view_placement(root_task_id: bookmark ref, range: [start, ))`
- view → view: check cache
- view → taskLink: run DFS, with child passing when: degree + passed_nodes < starting_index and terminate growing view when: passed_nodes > ing_index
- view → owner: placement array for range [start, )
- owner → push: create connection for future incoming data

## Create task

- owner → task: `create_task(user_id, root_task_id, super_task_id, task_details)`
- task → privilege: GET `roles(user_id, super_task_id)`
- task → task: CREATE task
- task → taskLink: CREATE link between `task_id`, `super_task_id`
- task → taskLink: COMPUTE degree, depth, index of task_id
- task → taskLink: RECOMPUTE (RECURSIVE to super-task) degree, depth, index of super_task_id

## Recomputing properties of a task link

- taskLink → taskLink: recompute
- taskLink → push: INVALIDATE BROWSER CACHE for task_id
- taskLink → push: INVALIDATE BROWSER CACHE for placement_array

## Invite user as "maintainer" over "task_id"

- owner → privilege: CREATE invitation()
- privilege → privilege: check roles (authorization)
- privilege → push: push event

## Accept invitation

- push → collab: notification: "new invitation"
- collab → privilege: ds

## Unassign user as "maintainer"

- owner → privilege: ds

## task details

- owner → taskLink: GET degree, depth (for user)
- owner → history: GET history (inc. subtasks optional)
- owner → privilege: GET roles (inc. inherited)
