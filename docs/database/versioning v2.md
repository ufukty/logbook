## **Logbook with Version Controlling v2**

**Objectives**

| **Objective**                        | **Outcome**                                                                                     |
| ------------------------------------ | ----------------------------------------------------------------------------------------------- |
| Fast forward not-conflicting changes |                                                                                                 |
| Checkout on super checkouts sub      | Parent objective should know which version of sub object it points to in each version of itself |
| Updates on sub iterates super        | Parents _subscribe_ to updates on sub objectives. (Happened after switch)                       |

### DB

_\*Tables are simplified_

### Procedures

$O_i$ is Object, $V_i$ is Version, $C_i$ is Commit, $A_i$ is action

#### Update $O_i$ ($V_i$) with $C_i$

```ruby
IterateParents(Obj, Vi, Id) {
	for l ∈ Links[Obj]
    if l.Sub.Ver == Vi
      # dsds
      l.Sup = Id
}

# Apply changes on Oi
nextVid = ApplyAction(Oi, Vi, Ci)
iterateParents(Oi, Vi, nextVid)
```

#### Linearize Subtree of $O_i$ on $V_i$

#### Revert $O_i$ to $V_{i - 1}$

$O_i$.Super.Links.Replace($O_i$, $O_i$.VersionAt($V_i$))

#### Create proposal on $O_i$ with $C_i = \{ A_1, A_2, A_3, \cdots \}$
