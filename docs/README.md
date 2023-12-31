# Logbook

## Versioning

**Pros**

-   Shorter user training
-   Auto fix update conflicts
-   Change proposals

**Cons**

-   Premature optimization for the database
-   Will postpone the release
-   Hard to detect bugs
-   Expensive to run
-   Ambiguity of scaling (upper objectives will have get updated more frequently than below ones)

**Challanges**

-   Integration with collaboration, requirements features. (due to storing computed props.)
-   Versioned chat may lead to have too many versions of super-objectives to manage, to scale

### Implementations

#### Enabled for every objective

-   Capable of saving multiple user-initiated actions on same resource without overhead of unnecessary duplication of unchanged resources (kind of _incremental backup_).

| Objective                            | Outcome                                                                                         |
| ------------------------------------ | ----------------------------------------------------------------------------------------------- |
| Fast forward not-conflicting changes |                                                                                                 |
| Checkout on super checkouts sub      | Parent objective should know which version of sub object it points to in each version of itself |
| Updates on sub iterates super        | Parents _subscribe_ to updates on sub objectives. (Happened after switch)                       |

##### Database

##### Procedures

$O_i$ is Object, $V_i$ is Version, $C_i$ is Commit, $A_i$ is action

-   Update $O_i$ ($V_i$) with $C_i$
-   Linearize Subtree of $O_i$ on $V_i$
-   Revert $O_i$ to $V_{i - 1}$
-   Create proposal on $O_i$ with $C_i = \{ A_1, A_2, A_3, \cdots \}$

#### When versions are tracked only for the subtree of selected objective

## Database Migration

```sh
cd database-migration
psql -d postgres -U ufuktan -f create.sql
```
