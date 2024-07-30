# Logbook database docs

## Document model

### Versioned documents

![](branching.png)

### Versioned trees

Design constraints

- Balance the time complexity and storage need between: 
  - **write**: saving changes into database for the next version, 
  - **read**: time to build a part of the document at any version
- Taking full snapshots of document is impossible (So, event-sourcing)
- Collaborators should be able to download only the changed parts rather than fully download or upload whole document at each sending/receiving change.

Design

- A content change in an objective gets saved as a new row in `objective` table with a new `vid` and same `oid`.
- Starting from the first parent of the updated objective, each ancestor gets a new row in the table with same `oid` but new `vid`.
- Objectives and links are needed to be immutable. Because of when they are siblings of updated objectives they will be referred from different versions of the document in `link` table.

Pros:

- Balanced read/writes: 
  - Storage effective
  - Checkouts are fast (compared to event-sourcing)

Cons:

- Implementation complexity.
- Calculated props are only stored for the active version.
- View builder needs to know which versions are ahead of the active version after checking out to previous version, which might not scale well for big version history.

![](simplified.png)

## View builder

Definition

- View builder is the linearization component for the hierarchically stored objectives, that generates an array of objectives that sit in the part of a document shown in the viewport of browser.

Design enablers

- Building the whole of document is almost never required.

Constraints:

- View builder needs the size of subtree of every node to be known for each version.

Design:

-	View builder should start comparing rows with only the latest version number 
-	As the parent node's version can never be older than its children's, view builder can pop the version numbers from version array when they are ahead of the held one, as it dig deeper. So, it can also postpone asking for the whole of version array, as it can fetch previous versions as it dig deeper and see another version than it holds
-	View builder may leverage `created_at` columns to quickly filter out ahead rows.



## Distributed transactions

![](transactions.png)
