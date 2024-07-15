# Roles and rights

## Table

| Action \ Role                         | Collaboration<br>Manager | Collaborator                                    | Viewer |
| ------------------------------------- | ------------------------ | ----------------------------------------------- | ------ |
| See tasks and task details            | ✔︎                        | ✔︎                                               | ✔︎      |
| Extend tasks with subtasks            | ✔︎                        | ✔︎                                               |        |
| Edit task content                     | ✔︎                        | ✔︎ (created tasks & and any-level subtasks)      |        |
| Change task ordering                  | ✔︎                        | ✔︎ (any-level subtasks of created tasks)         |        |
| Subscribe completion progress changes | ✔︎                        | ✔︎                                               |        |
| Subscribe task detail changes         | ✔︎                        | ✔︎                                               |        |
| Export blueprint                      | ✔︎                        | ✔︎                                               |        |
| Assign / invite users                 | ✔︎                        | ✔︎ (except users directly restricted by manager) |        |
| Restrict users                        | ✔︎                        | ✔︎ (except users directly invited by manager)    |        |

## Notes

- Collaborators can create subtasks and become **creator of subtasks** they create. 

- Collaboration role inherits as itself to sub-tasks.

- Creator role inherits as collaborator.

- Example for role inhertance. Let's say:

  - **Alice** has created a task named **First Task** then invited **Bob** to collaborate on it. **Bob** has accepted invitation.

    ```
    - First Task
    ```

  - Alice has created a subtask for **First Task** which she named as **Second Task**. 
    ```
    - First Task
    	- Second Task <--
    ```

  - Then, **Bob** has added another step for **First Task** as **Third Task**. 
    ```
    - First Task
    	- Second Task
    	- Third Task <--
    ```

  - Roles are inherited by default from super-tasks to sub-tasks. So, here is the final state of role sharing for each task:

    ```
    - First Task    (Alice is creator,      Bob is collaborator) <--
      - Second Task (Alice is creator,      Bob is collaborator) <-- inherited
      - Third Task  (Alice is collaborator, Bob is creator)      <-- Bob is creator, not Alice.
    ```

- Creators or inherited-creators can not degranted from gaining inherited-creator role on created subtasks. If this is really needed, creator of subtask can select "solo subtask" option.

- Each subtask inherits roles applied to *primary supertask* by default as long as not overwritten at its level. Including additional users on subtasks or excluding users from supertasks are possible.

- After a collaborator "unassigned" from previously assigned task, and loose its connection with subtasks its created:

  - In-collaboration subtasks of unassigned user:
    - Those tasks meant to be presented to creators of supertasks. 
    - Creators can "adopt" those tasks or delete them from tree.
  - Private sub-tasks of unassigned user switch:
    - Those tasks are switched to "detached" mode.
    - Detached tasks are only reachable for creator, from bookmarks that point to root of each tree.
    - Undelivered subtasks loose their impact on completion progress of super-task.
