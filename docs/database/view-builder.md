# View builder

Definition

- View builder is the linearization component for the hierarchically stored objectives, that generates an array of objectives that sit in the part of a document shown in the viewport of browser.

Design enablers

- Building the whole of document is almost never required.

Constraints:

- View builder needs the size of subtree of every node to be known for each version.

Design:

- View builder should start comparing rows with only the latest version number
- As the parent node's version can never be older than its children's, view builder can pop the version numbers from version array when they are ahead of the held one, as it dig deeper. So, it can also postpone asking for the whole of version array, as it can fetch previous versions as it dig deeper and see another version than it holds
- View builder may leverage `created_at` columns to quickly filter out ahead rows.
